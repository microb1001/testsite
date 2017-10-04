package main

import (
"html/template"
"log"
"net/http"
"strings"
"time"
"./csv"
"fmt"
"strconv"
"sort"
"./webelements"
)

type good struct {
	UIN int `csv:"UIN"`
	Barcode int
	VendorCode string `csv:"Артикул"`
	Brief string `csv:"Описание"`
	Price int `csv:"Цена"`
	Quantity int `csv:"Количество"`
	Available bool `csv:"В продаже"`
	MainCategory string `csv:"Категория"`
	Category string `csv:"Товар"`
	Spec map[string]string `csv:"Поиск"`
	Pictures string   `csv:"Артикул"`
	Info  int  `csv:"N"`
	ShortDescription string `csv:"Описание"`
	Description string   `csv:"Характеристика"`
	Images string
	UrlAlias string `csv:"Path"`
}

type cart_type struct {
	VendorCode string
	Quantity int
}

type user struct {
	UID string
	login string
	passhash string
	email string
	cookie string
	cart []cart_type
	shipping_adress1 string
	shipping_adress2 string
	shipping_adress3 string
	payment_info string
}

const items_per_page=2
var goods []good
var goodsmap map[string]int = make(map[string]int, len(goods))
var sel map[string][]int = make(map[string][]int,50)

type e struct{
	name,url string
	elem []string
}
var category1e []e = []e{
	e{"/stolovoe-serebro","Столовое серебро",[]string{}},
	e{"/ukrashenia","Украшения",[]string{}},
	e{"/zoloto","Золото",[]string{}},
	e{"/prochee","Прочее",[]string{}},
	}
var category1 map[string]map[string]string = make(map[string]map[string]string,50)
var category2 map[string]string = make(map[string]string,50)
type  category1listType struct{	Key string; Value [] struct{ Key, Url string}}
var category1list [] category1listType
var category2list [] struct{key, value string}

func main() {
	mycsv.Load_csv(&goods,"list.csv", "csv")
	sel["/"] = []int{0,1,2,3,4,5,6,7,8,9,10,}
	for i,k:=range goods {
	goodsmap[k.VendorCode]=i
	sel[k.UrlAlias]=append(sel[k.UrlAlias],i)

	if category1[k.MainCategory]==nil {
		category1[k.MainCategory]=make(map[string]string,50)
	}
	category1[k.MainCategory][k.Category]=k.UrlAlias

	if category2[k.Category]=="" {
		category2[k.Category]=k.UrlAlias
	}
	}
	for f,g:= range category1 {
		var tp category1listType
		tp.Key =f
		for k,m:= range g {
			tp.Value =append(tp.Value,struct{Key, Url string}{k,m})
		}

		category1list=append(category1list,tp)
		sort.Slice(tp.Value, func(i, j int) bool { return tp.Value[i].Key < tp.Value[j].Key })
	}
	sort.Slice(category1list, func(i, j int) bool { return len(category1list[i].Value) > len(category1list[j].Value) }) // по количеству товаров

	for k,m:= range category2 {
		category2list=append(category2list,struct{key, value string}{k,m})

	}
	sort.Slice(category2list, func(i, j int) bool { return category2list[i].key < category2list[j].key })

	//fmt.Println(category1,"----",category2,"tt ")
	fmt.Println(category1list)
	fmt.Println(category2list)
	//mycsv.Dump(goods)

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/product/", imageHandler)
	fs := http.FileServer(http.Dir("img/"))
	http.Handle("/images/", http.StripPrefix("/images/", fs)) // небезопасно отдает файлы любого типа!
	log.Fatal(http.ListenAndServe("localhost:80", nil))
}



// The default template includes two template blocks ("sidebar" and "content")
// that may be replaced in templates derived from this one.
var mainTemplate = template.Must(template.ParseFiles("index.tmpl"))

// mainHandler is an HTTP handler that serves the index page (list of goods).
func mainHandler(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	type LinkType struct {
		good
		URL, Title string
	}

	var data struct{
		Title, Body string
		Links []LinkType
		Pager webelements.PagerType
		Cat [] category1listType
		Timer time.Duration //Timer
	}

	//_=r.ParseForm() // Само вызывается из FormValue
	fmt.Println("URL.Path: ",r.URL.Path," RawPath: ",r.URL.RawPath," RequestURI():",r.URL.RequestURI(),"Host: ",r.Host,"FormValue: ",r.FormValue("p"))

	var cnt int = 0
	var pageCurrent int
	var mainPage string

	data.Links =make([]LinkType,0,items_per_page)
	data.Title= "Image gallery 11-11"
	data.Body = "Welcome to the image gallery."
	mainPage=r.URL.Path // !!Нужна обработка пользовательского ввода!!
	pageCurrent,err:=strconv.Atoi(r.FormValue("p"))
	if err != nil {
		pageCurrent = 0
	}
	pageMax :=(len(sel[mainPage])-1)/items_per_page // начинается с нуля
	data.Pager=webelements.Pager(0, pageMax,pageCurrent,mainPage)
	for _,i := range sel[mainPage][webelements.MinMax((pageCurrent)*items_per_page,0,len(sel[mainPage])):webelements.MinMax((pageCurrent+1)*items_per_page,0,len(sel[mainPage]))] {
		data.Links = append(data.Links, LinkType{goods[i], "/product/" + goods[i].VendorCode, ""})
	}

	data.Cat=category1list
	data.Timer=time.Now().Sub(start)
	if err := mainTemplate.Execute(w, data); err != nil {
		log.Println(err)
	}
	t := time.Now()
	elapsed := t.Sub(start)
	log.Println("timer ",cnt, elapsed)
}

// imageTempl9ate is a clone of indexTemp9late that provides
// alternate "sidebar" and "content" templates.
var imageTemplate = template.Must(template.Must(mainTemplate.Clone()).ParseFiles("image.tmpl"))

// Image is a data structure used to populate an imageTemp9late.
type Image struct {
	Title string
	URL   string
}

// imageHandler is an HTTP handler that serves the image pages.
func imageHandler(w http.ResponseWriter, r *http.Request) {

	type spec1type struct {
		Key, Value string
	}

	type datatype struct {
		good
		Spec1 []spec1type
		Title string
		URL   string
	}

	dataindex, ok := goodsmap[strings.TrimPrefix(r.URL.Path, "/product/")]
	//data, ok := images[strings.TrimPrefix(r.URL.Path, "/product/")]
	if !ok {
		http.NotFound(w, r)
		return
	}
	var data datatype=datatype{goods[dataindex],
	[]spec1type{},"","/images/"+goods[dataindex].VendorCode+".jpg"}
	for _,item3 := range []string{"Высота","Ширина","Диаметр","Размер"} {
	if data.Spec[item3]!=""{data.Spec1=append(data.Spec1, spec1type{item3,data.Spec[item3]})

	}
	}

	if err := imageTemplate.Execute(w, data); err != nil {
		log.Println(err)
	}
}

