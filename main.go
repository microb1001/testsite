package main

import (
"html/template"
"log"
"net/http"
"strings"
"time"
"fmt"
"strconv"
"./webelements"
"./mydb"
)

const items_per_page=2
var goods mydb.Goods





func main() {
	goods.Init("list.csv")
	goods.Sel["/"] = []int{0,1,2,3,4,5,6,7,8,9,10,}

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
		mydb.Good
		URL, Title string
	}

	var data struct{
		Title, Body string
		Links []LinkType
		Pager webelements.PagerType
		Cat []mydb.Category1listType
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
	var l, h int
	data.Pager, l, h = webelements.Pager(pageCurrent,items_per_page, len(goods.Sel[mainPage]), mainPage+"?")
	for _,i := range goods.Sel[mainPage][l:h] {
		data.Links = append(data.Links, LinkType{goods.O[i], "/product/" + goods.O[i].VendorCode, ""})
	}

	data.Cat=goods.Category1list
	data.Timer=time.Now().Sub(start)
	fmt.Println(data.Timer.String())
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
		mydb.Good
		Spec1 []spec1type
		Title string
		URL   string
	}

	dataindex, ok := goods.Goodsmap[strings.TrimPrefix(r.URL.Path, "/product/")]
	//data, ok := images[strings.TrimPrefix(r.URL.Path, "/product/")]
	if !ok {
		http.NotFound(w, r)
		return
	}
	var data datatype=datatype{goods.O[dataindex],
	[]spec1type{},"","/images/"+goods.O[dataindex].VendorCode+".jpg"}
	for _,item3 := range []string{"Высота","Ширина","Диаметр","Размер"} {
	if data.Spec[item3]!=""{data.Spec1=append(data.Spec1, spec1type{item3,data.Spec[item3]})

	}
	}

	if err := imageTemplate.Execute(w, data); err != nil {
		log.Println(err)
	}
}

