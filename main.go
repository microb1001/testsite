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
	"index/suffixarray"
)

const items_per_page=10
var goods mydb.Goods_type
var userCart map[uint64]mydb.Usercart_type
var index *suffixarray.Index
var Sphi webelements.Sphinx
//var Context webelements.SessionListType =

func main() {
	goods.Init("list.csv")
	goods.AddPrice("")
	//mycsv.Dump(goods)
	Sphi.Init()
	Sphi.Add(&goods)

	userCart = make(map[uint64]mydb.Usercart_type)
	var tmpstring []byte
	for _,k:=range goods.O {
		tmpstring=append(tmpstring,[]byte(k.Description)...)
	}
	fmt.Println(tmpstring)
	index = suffixarray.New(tmpstring)
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/product/", imageHandler)
	http.HandleFunc("/cart/", cartHandler)
	http.HandleFunc("/search1/", searchHandler)
	fs1 := webelements.MyFs{http.Dir("img/")}
	//http.ListenAndServe(":8080", http.FileServer(fs1))
	//fs := http.FileServer(http.Dir("img/"))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer( fs1))) // небезопасно отдает файлы любого типа!
	log.Fatal(http.ListenAndServe("localhost:80", nil))
}

// The default template includes two template blocks ("sidebar" and "content")
// that may be replaced in templates derived from this one.
var mainTemplate = template.Must(template.ParseFiles("index.tmpl"))
// image Template is a clone of index Template that provides
// alternate "sidebar" and "content" templates.
var imageTemplate = template.Must(template.Must(mainTemplate.Clone()).ParseFiles("image.tmpl"))

var cartTemplate = template.Must(template.Must(mainTemplate.Clone()).ParseFiles("cart.tmpl"))
var searchTemplate = template.Must(template.Must(mainTemplate.Clone()).ParseFiles("search.tmpl"))
// mainHandler is an HTTP handler that serves the index page (list of goods).

func mainHandler(w http.ResponseWriter, r *http.Request) {

	type LinkType struct {
		mydb.GoodsElem_type

		URL, Title, Image string
	}
	var data struct{
		Title, Body string
		Links []LinkType
		Pager webelements.PagerType
		Cat []mydb.Category1List_type
		Timer time.Duration //Timer
		Session uint64
	}
	start := time.Now()
	data.Session =webelements.SessionGet(w,r)
	goods.Mu.RLock()
	defer goods.Mu.RUnlock()

	searchstring := r.FormValue("text")
	if searchstring != "" {
		otv:=Sphi.Find(searchstring)
		for _,i:=range otv{
			fmt.Println(goods.O[i].Description)
		}

		//SearchResult:=index.Lookup([]byte(searchstring), -1)
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
	if mainPage=="/search/" {
		goods.Sel[mainPage]=Sphi.Find(searchstring)
	}
	fmt.Println("<<<<<<<<<<<<<<<",goods.Sel[mainPage])
	pageCurrent,err:=strconv.Atoi(r.FormValue("p"))
	if err != nil {
		pageCurrent = 0
	}
	var l, h int
	data.Pager, l, h = webelements.Pager(pageCurrent,items_per_page, len(goods.Sel[mainPage]), mainPage+"?")
	for _,i := range goods.Sel[mainPage][l:h] {
		data.Links = append(data.Links, LinkType{goods.O[i], "/product/" + goods.O[i].VendorCode, "","/images/400/"+ goods.O[i].VendorCode+".jpg"})
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
		mydb.GoodsElem_type
		Spec1   []spec1type
		Title   string
		URL     string
		Session uint64
	}

	dataindex, ok := goods.Goodsmap[strings.TrimPrefix(r.URL.Path, "/product/")]
	//data, ok := images[strings.TrimPrefix(r.URL.Path, "/product/")]
	if !ok {
		http.NotFound(w, r)
		return
	}

	var data datatype = datatype{goods.O[dataindex],
		[]spec1type{}, "", "/images/" + goods.O[dataindex].VendorCode + ".jpg", webelements.SessionGet(w, r)}
	for _, item3 := range []string{"Высота", "Ширина", "Диаметр", "Размер"} {
		if data.Spec[item3] != "" {
			data.Spec1 = append(data.Spec1, spec1type{item3, data.Spec[item3]})

		}
	}

	if err := imageTemplate.Execute(w, data); err != nil {
		log.Println(err)
	}
}
func cartHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title, Body string
		mydb.User_type
		UserCart    mydb.Usercart_type
		Session     uint64
	}
	sessid := webelements.SessionGet(w, r)
	data.Session = sessid
	if r.FormValue("additem") != "" {
		goodsid := goods.Goodsmap[r.FormValue("additem")]
		userCart[sessid] = append(userCart[sessid], &goods.O[goodsid]) // refresh добавляет повтор убрать (можно через реферрер
	}
	data.UserCart = userCart[sessid]
	if err := cartTemplate.Execute(w, data); err != nil {
		log.Println(err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title, Body string
		mydb.User_type
		UserCart    mydb.Usercart_type
		Session     uint64
		SearchResult []int
	}
	sessid := webelements.SessionGet(w, r)
	data.Session = sessid
	searchstring := r.FormValue("text")
	fmt.Println("=============",[]byte(searchstring))
	if searchstring != "" {
		otv:=Sphi.Find(searchstring)
		fmt.Println(otv)
		for _,i:=range otv{
			fmt.Println(goods.O[i].Description)
		}

		data.SearchResult=index.Lookup([]byte(searchstring), -1)
	}

	data.UserCart = userCart[sessid]
	if err := searchTemplate.Execute(w, data); err != nil {
		log.Println(err)
	}
}