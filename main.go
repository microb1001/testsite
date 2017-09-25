package main

import (
"html/template"
"log"
"net/http"
"strings"
"time"

)

func main() {
	//parse2("list.csv",[]string{"N","UIN", "Артикул", "Описание", "Количество", "Вес", "Цена"})
	goods=append(goods, good{Info:"1231"})
	goods=append(goods, good{Articul:"1231"})
	parse(&goods,"list.csv")
	goods=append(goods, good{Articul:"5555"})
	dump(goods,"list.csv")

	items[550]=1
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/image/", imageHandler)
	log.Fatal(http.ListenAndServe("localhost:80", nil))
}

var items [100000]int
var goods [] good
var lst []int
var cnt int = 0

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

type good struct {
	Articul string   `json:"N"`
	Info  string   `json:"UIN"`
	Image string `json:"Артикул"`

}

// indexHandler is an HTTP handler that serves the index page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	lst=make([]int,0,120000)
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
	for item := range items {
		if item==1 {
			cnt=cnt+1
		}
		lst=append(lst,item)
	}
	if err := indexTemplate.Execute(w, data); err != nil {
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
