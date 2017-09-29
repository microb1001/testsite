package main

import (

"html/template"
"log"
"net/http"
"strings"
"time"
"./csv"
	"fmt"
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
	Category []string `csv:"Товар"`
	Pictures string   `csv:"Артикул"`
	Info  int  `csv:"N"`
	ShortDescription string `csv:"Описание"`
	Description string   `csv:"Характеристика"`
	Images string
}

var items [100000]good
var goods [] good
var sel []int

func main() {
	//sel = []int{1, 2, 3,4,200,280,600,860,5,1100,444,555,556,667,668,669,4,6,8,888}
	sel = []int{0,1,2,3,4,5,6,7,8,9,10}
	mycsv.Load_csv(&goods,"list.csv", "csv")
	//mycsv.Dump(goods)
	for i:=range sel {
		items[i].VendorCode="Q1"
	}
	items[550].VendorCode="Q1"
	items[1551].UIN=1


	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/image/", imageHandler)
	log.Fatal(http.ListenAndServe("localhost:80", nil))
}



// indexTemplate is the main site template.
// The default template includes two template blocks ("sidebar" and "content")
// that may be replaced in templates derived from this one.
var indexTemplate = template.Must(template.ParseFiles("index.tmpl"))

// Index is a data structure used to populate an indexTemplate.
type Index struct {
	Title string
	Body  string
	Links []Link
}

type Link struct {
	URL, Title string
}


// indexHandler is an HTTP handler that serves the index page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	var data2 struct{
		Links []good
		Title, Body string
	}
	var cnt int = 0
	data2.Links =make([]good,0,120000)
	data2.Title= "Image gallery 11-11"
	data2.Body = "Welcome to the image gallery."
	start := time.Now()

	data := &Index{
		Title: "Image gallery 11-11",
		Body:  "Welcome to the image gallery.",
	}

	for name, img := range images {
		data.Links = append(data.Links, Link{
			URL:   "/image/" + name,
			Title: img.Title,
		})
	}
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
	for _,i := range sel {
		data2.Links =append(data2.Links,goods[i])
	}
	if err := indexTemplate.Execute(w, data2); err != nil {
		log.Println(err)
	}
	t := time.Now()
	elapsed := t.Sub(start)
	log.Println("timer ",cnt, elapsed)
}

// imageTemplate is a clone of indexTemplate that provides
// alternate "sidebar" and "content" templates.
var imageTemplate = template.Must(template.Must(indexTemplate.Clone()).ParseFiles("image.tmpl"))

// Image is a data structure used to populate an imageTemplate.
type Image struct {
	Title string
	URL   string
}

// imageHandler is an HTTP handler that serves the image pages.
func imageHandler(w http.ResponseWriter, r *http.Request) {
	data, ok := images[strings.TrimPrefix(r.URL.Path, "/image/")]
	if !ok {
		http.NotFound(w, r)
		return
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
