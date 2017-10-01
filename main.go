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
	Category []string `csv:"Товар"`
	Pictures string   `csv:"Артикул"`
	Info  int  `csv:"N"`
	ShortDescription string `csv:"Описание"`
	Description string   `csv:"Характеристика"`
	Images string
	UrlAlias string
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
var items [100000]good
var goods [] good
var goodsmap map[string]int = make(map[string]int, len(goods))
var sel []int




func main() {
	//sel = []int{1, 2, 3,4,200,280,600,860,5,1100,444,555,556,667,668,669,4,6,8,888}
	sel = []int{0,1,2,3,4,5,6,7,8,9,10,}
	mycsv.Load_csv(&goods,"list.csv", "csv")
	for i,k:=range goods {
	goodsmap[k.VendorCode]=i
	}
	//mycsv.Dump(goods)
	for i:=range sel {
		items[i].VendorCode="Q1"
	}
	items[550].VendorCode="Q1"
	items[1551].UIN=1


	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/product/", imageHandler)
	fs := http.FileServer(http.Dir("img/"))
	http.Handle("/images/", http.StripPrefix("/images/", fs)) // небезопасно отдает файлы любого типа!
	// This works too, but "/static2/" fragment remains and need to be striped manually
	//http.HandleFunc("/static2/", func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, r.URL.Path[1:])
	// })
	//NOTE Serving up a filesystem naively is a potentially dangerous thing (there are potentially ways to break out of the rooted tree) hence I recommend that unless you really know what you're doing, use http.FileServer and http.Dir as they include checks to make sure people can't break out of the FS, which ServeFile doesn't.
/*  пример управляемой FileServer
package main

	import (
		"net/http"
	"os"
	)

	type justFilesFilesystem struct {
		fs http.FileSystem
	}

	func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
	return nil, err
	}
	return neuteredReaddirFile{f}, nil
	}

	type neuteredReaddirFile struct {
		http.File
	}

	func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
	}

	func main() {
		fs := justFilesFilesystem{http.Dir("/tmp/")}
		http.ListenAndServe(":8080", http.FileServer(fs))
	} */





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

	const items_per_page=2

	type PagerType struct{
		Page int
		Class string
		Url string
		Current bool
	}

	type Link2 struct {
		good
		URL, Title string

	}

	var data2 struct{
		Links []Link2
		Pager []PagerType
		Title, Body string
	}


	//_=r.ParseForm() // и так вызывается из FormValue
	fmt.Println("URL.Path: ",r.URL.Path," RawPath: ",r.URL.RawPath," RequestURI():",r.URL.RequestURI(),"Host: ",r.Host,"FormValue: ",r.FormValue("p"))
	var cnt int = 0
	data2.Links =make([]Link2,0,120000)
	data2.Title= "Image gallery 11-11"
	data2.Body = "Welcome to the image gallery."
	var ipage int;

	ipage,err:=strconv.Atoi(r.FormValue("p"))
	if err != nil {
		ipage=1
	}
	start := time.Now()

	data := &Index{
		Title: "Image gallery 11-11",
		Body:  "Welcome to the image gallery.",
	}
	maxi:=(len(sel)-1)/items_per_page+1
	for ii:=minMax(ipage-2,1,maxi);ii<=minMax(ipage+2,1,maxi);ii++{
		data2.Pager=append(data2.Pager,PagerType{ii,"","?p="+strconv.Itoa(ii),ii==ipage} )
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

	for _,i := range sel[minMax((ipage-1)*items_per_page,0,len(sel)):minMax(ipage*items_per_page,0,len(sel))] {
		data2.Links =append(data2.Links,Link2{goods[i],"/product/"+goods[i].VendorCode,""})
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
	type datatype struct {
		good
		Title string
		URL   string
	}

	dataindex, ok := goodsmap[strings.TrimPrefix(r.URL.Path, "/product/")]
	//data, ok := images[strings.TrimPrefix(r.URL.Path, "/product/")]
	if !ok {
		http.NotFound(w, r)
		return
	}
	var data datatype=datatype{goods[dataindex],"","/images/"+goods[dataindex].VendorCode+".jpg"}

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
