package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"fmt"
	"strconv"
	web "./webelements"
	db "./mydb"
	"./webelements/session"
	"./webelements/sphinx"
	"./webelements/pager"
	"./webelements/escape"
	//"./mydb/goods"
)

const items_per_page=10
var goods db.Goods_type
var userCart map[uint64]db.Usercart_type
var Sphi sphinx.Sphinx

func main() {
	goods.Init("list.csv")
	goods.AddPrice("")
	//mycsv.Dump(goods)
	Sphi.Init()
	Sphi.Add(&goods)
	fmt.Println(string(escape.Go([]byte("100asd5fgфывап.:/"), escape.MakeFn(append(escape.Eng,escape.Rus...)))))
	userCart = make(map[uint64]db.Usercart_type)
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/product/", imageHandler)
	http.HandleFunc("/cart/", cartHandler)
	fs1 := web.MyFs{http.Dir("img/")}
	//http.ListenAndServe(":8080", http.FileServer(fs1))
	//fs := http.FileServer(http.Dir("img/"))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer( fs1)))
	log.Fatal(http.ListenAndServe("localhost:80", nil))
}

// The default template includes two template blocks ("sidebar" and "content")
// that may be replaced in templates derived from this one.
var mainTemplate = template.Must(template.ParseFiles("index.tmpl","paginator.tmpl"))
// image Template is a clone of index Template that provides
// alternate "sidebar" and "content" templates.
var imageTemplate = template.Must(template.Must(mainTemplate.Clone()).ParseFiles("image.tmpl"))

var cartTemplate = template.Must(template.Must(mainTemplate.Clone()).ParseFiles("cart.tmpl"))
var searchTemplate = template.Must(template.Must(mainTemplate.Clone()).ParseFiles("search.tmpl"))

// mainHandler is an HTTP handler that serves the index page (list of goods).
func mainHandler(w http.ResponseWriter, r *http.Request) {

	type LinkType struct {
		db.GoodsElem_type

		URL, URLtoCart, Title, Image string
		Separator3 int
	}
	var data struct{
		Title, Body string
		Links       []LinkType
		Pager       pager.Pager
		Cat         []db.Category1List_type
		Timer       time.Duration //Timer
		Session     uint64
	}
	var p []string
	start := time.Now()
	data.Session = session.Get(w,r)
	goods.Mu.RLock()
	defer goods.Mu.RUnlock()

	searchstring := r.FormValue("text")
	if searchstring != "" {
		p=append(p,"text="+searchstring)
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
	pageCurrent,err:=strconv.Atoi(r.FormValue("p"))
	if err != nil {
		pageCurrent = 0
	}
	var l, h int
	mainPage2:=mainPage+"?"
	for _,i:=range p{
		mainPage2=mainPage2+"&"+i
	}
	data.Pager, l, h = pager.Set(pageCurrent,items_per_page, len(goods.Sel[mainPage]), mainPage2+"&")
	for t,i := range goods.Sel[mainPage][l:h] {
		data.Links = append(data.Links, LinkType{goods.O[i], "/product/" + goods.O[i].VendorCode, "/cart/?additem=" + goods.O[i].VendorCode,"","/images/400/"+ goods.O[i].VendorCode+".jpg",t%6})
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

// imageHandler is an HTTP handler that serves the image pages.
func imageHandler(w http.ResponseWriter, r *http.Request) {

	type spec1type struct {
		Key, Value string
	}

	type datatype struct {
		db.GoodsElem_type
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
		[]spec1type{}, "", "/images/" + goods.O[dataindex].VendorCode + ".jpg", session.Get(w, r)}
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
		db.User_type
		UserCart    db.Usercart_type
		Session     uint64
	}
	sessid := session.Get(w, r)
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
