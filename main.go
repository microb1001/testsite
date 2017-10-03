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

type PagerType struct{
	Page int
	Class string
	Url string
	Current bool
}

const items_per_page=2
var goods []good
var goodsmap map[string]int = make(map[string]int, len(goods))
var sel map[string][]int = make(map[string][]int,50)

func main() {
	mycsv.Load_csv(&goods,"list.csv", "csv")
	sel["/"] = []int{0,1,2,3,4,5,6,7,8,9,10,}
	for i,k:=range goods {
	goodsmap[k.VendorCode]=i
	sel[k.UrlAlias]=append(sel[k.UrlAlias],i)
	}
	//mycsv.Dump(goods)

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/product/", imageHandler)
	fs := http.FileServer(http.Dir("img/"))
	http.Handle("/images/", http.StripPrefix("/images/", fs)) // небезопасно отдает файлы любого типа!
	log.Fatal(http.ListenAndServe("localhost:80", nil))
}

func minMax(index,min,max int) int{
	if index<min {
		return min
	}
	if index>max {
		return max
	}
	return index
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
		Pager []PagerType
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
	mainPage = "/kuvshin"
	pageCurrent,err:=strconv.Atoi(r.FormValue("p"))
	if err != nil {
		pageCurrent = 0
	}
	pageMax :=(len(sel[mainPage])-1)/items_per_page // начинается с нуля
	for ii:=minMax(pageCurrent-2,0, pageMax);ii<=minMax(pageCurrent+2,0, pageMax);ii++{
		data.Pager=append(data.Pager,PagerType{ii+1,"","?p="+strconv.Itoa(ii),ii == pageCurrent} )
	}

//	for name, img := range images {
//		data.Links = append(data.Links, Link{
//			URL:   "/image/" + name,
//			Title: img.Title,
//		})
//	}
	/*for _,item := range items {
		if item.UIN==1 {
			cnt=cnt+1
			lst=append(lst,item)
		}
		if item.Articul=="Q1" {
			cnt=cnt+1
			lst=append(lst,item)
		}
		//lst=append(lst,item)

		//
	} */

	data.Timer=time.Now().Sub(start)
	for _,i := range sel[mainPage][minMax((pageCurrent)*items_per_page,0,len(sel[mainPage])):minMax((pageCurrent+1)*items_per_page,0,len(sel[mainPage]))] {
			data.Links =append(data.Links, LinkType{goods[i],"/product/"+goods[i].VendorCode,""})
	}
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

// images specifies the site content: a collection of images.   https://lenta.ru/
var images = map[string]*Image{
	"go":     {"The Go Gopher", "https://golang.org/doc/gopher/frontpage.png"},
	"google": {"The Google Logo", "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"},
	"lenta": {"The Lenta Logo", "https://icdn.lenta.ru/images/2017/09/20/18/20170920184345749/pic_9236d1c3d84b722a85fe66166a0ef251.jpg"},
}
