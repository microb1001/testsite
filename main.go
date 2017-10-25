package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"strconv"
	"./webelements/session"
	"./webelements/sphinx"
	"./webelements/pager"
	_"./webelements/escape"
	"./webelements/myfs"
	"./mydb/good7"
	"./mydb/cart7"
	"./templates/tdata"
)

const items_per_page=10
var goods good7.Goods
var userCart map[uint64]cart7.Cart
var sphi sphinx.Sphinx

func main() {
	goods.Init("list.csv")
	goods.AddPrice("")
	sphi.Init()
	sphi.Add(&goods)
	//mycsv.Dump(goods)

	userCart = make(map[uint64]cart7.Cart)

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/product/", imageHandler)
	http.HandleFunc("/cart/", cartHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(
												myfs.FileSystem{http.Dir("img/")})))
	log.Fatal(http.ListenAndServe("localhost:80", nil))
}

var mainTemplate = template.Must(template.ParseFiles("templates/index.tmpl","templates/paginator.tmpl"))
var imageTemplate = template.Must(template.Must(mainTemplate.Clone()).ParseFiles("templates/image.tmpl"))
var cartTemplate = template.Must(template.Must(mainTemplate.Clone()).ParseFiles("templates/cart.tmpl"))

// mainHandler is an HTTP handler that serves the index page (list of goods).
func mainHandler(w http.ResponseWriter, r *http.Request) {

	var data tdata.List
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
	//fmt.Println("URL.Path: ",r.URL.Path," RawPath: ",r.URL.RawPath," RequestURI():",r.URL.RequestURI(),"Host: ",r.Host,"FormValue: ",r.FormValue("p"))

	var cnt int = 0
	var pageCurrent int
	var mainPage string

	data.Links =make([]tdata.LinkType,0,items_per_page)
	data.Title= "Image gallery 11-11"
	data.Body = "Welcome to the image gallery."
	mainPage=r.URL.Path // !!Нужна обработка пользовательского ввода!!
	if mainPage=="/search/" {
		goods.Sel[mainPage]= sphi.Find(searchstring)
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
		data.Links = append(data.Links, tdata.LinkType{goods.O[i], "/product/" + goods.O[i].VendorCode, "/cart/?additem=" + goods.O[i].VendorCode,"","/images/400/"+ goods.O[i].VendorCode+".jpg",t%6})
	}

	data.Cat=goods.Category1list
	data.Timer=time.Now().Sub(start)
	if err := mainTemplate.Execute(w, data); err != nil {
		log.Println(err)
	}
	t := time.Now()
	elapsed := t.Sub(start)
	log.Println("timer ",cnt, elapsed)
}

// imageHandler is an HTTP handler that serves the info about good.
func imageHandler(w http.ResponseWriter, r *http.Request) {


	dataindex, ok := goods.Goodsmap[strings.TrimPrefix(r.URL.Path, "/product/")]
	//data, ok := images[strings.TrimPrefix(r.URL.Path, "/product/")]
	if !ok {
		http.NotFound(w, r)
		return
	}

	var data tdata.Product = tdata.Product{goods.O[dataindex],
		[]tdata.ProductSpec{}, "", "/images/" + goods.O[dataindex].VendorCode + ".jpg", session.Get(w, r)}
	for _, item3 := range []string{"Высота", "Ширина", "Диаметр", "Размер"} {
		if data.Spec[item3] != "" {
			data.Spec1 = append(data.Spec1, tdata.ProductSpec{item3, data.Spec[item3]})
		}
	}

	if err := imageTemplate.Execute(w, data); err != nil {
		log.Println(err)
	}
}

// отображает корзину
func cartHandler(w http.ResponseWriter, r *http.Request) {
	var data tdata.Cart
	sessid := session.Get(w, r)
	data.Session = sessid
	if r.FormValue("additem") != "" {
		goodsid := goods.Goodsmap[r.FormValue("additem")]
		userCart[sessid] = append(userCart[sessid], cart7.Elem{"+",1,goods.O[goodsid]}) // refresh добавляет повтор убрать (можно через реферрер
	}
	data.UserCart = userCart[sessid]
	data.TotalPrice = 0
	data.TotalCount =0
	for _,m:= range userCart[sessid]{
	 data.TotalPrice+=m.Price
	 data.TotalCount+=1
	}
	if err := cartTemplate.Execute(w, data); err != nil {
		log.Println(err)
	}
}
