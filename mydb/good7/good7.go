package good7

import (
	"sync"
	"sort"
	"../../csv"
)

type Elem struct {
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

type  Category1List_type struct{Key string; Value [] struct{ Key, Url string}}

type Goods struct {
	Mu            sync.RWMutex
	O             []Elem

	Goodsmap      map[string]int
	Sel           map[string][]int
	category1     map[string]map[string]string
	category2     map[string]string
	Category1list [] Category1List_type
	category2list [] struct{key, value string}
}
func (s *Goods) Init (filename string){
	s.Mu.Lock()
	defer s.Mu.Unlock()
	mycsv.Load_csv(&s.O, filename, "csv")
	s.Goodsmap = make(map[string]int, len(s.O))
	s.Sel= make(map[string][]int,50)
	s.category1 = make(map[string]map[string]string,50)
	s.category2 = make(map[string]string,50)
	s.Sel["/"] = []int{0,1,2,3,4,5,6,7,8,9,10,}
	for i,k:=range s.O {
		s.Goodsmap[k.VendorCode]=i
		s.Sel[k.UrlAlias]=append(s.Sel[k.UrlAlias],i)

		if s.category1[k.MainCategory]==nil {
			s.category1[k.MainCategory]=make(map[string]string,50)
		}
		s.category1[k.MainCategory][k.Category]=k.UrlAlias

		if s.category2[k.Category]=="" {
			s.category2[k.Category]=k.UrlAlias
		}
	}
	for f,g:= range s.category1 {
		var tp Category1List_type
		tp.Key =f
		for k,m:= range g {
			tp.Value =append(tp.Value,struct{Key, Url string}{k,m})
		}

		s.Category1list =append(s.Category1list,tp)
		sort.Slice(tp.Value, func(i, j int) bool { return tp.Value[i].Key < tp.Value[j].Key })
	}
	sort.Slice(s.Category1list, func(i, j int) bool { return len(s.Category1list[i].Value) > len(s.Category1list[j].Value) }) // по количеству товаров

	for k,m:= range s.category2 {
		s.category2list=append(s.category2list,struct{key, value string}{k,m})

	}
	sort.Slice(s.category2list, func(i, j int) bool { return s.category2list[i].key < s.category2list[j].key })
}

func (s *Goods) AddPrice (filename string){
	s.Mu.Lock()
	defer s.Mu.Unlock()
}

func (s *Goods) ToByte (i int) []byte {
	return []byte(s.O[i].Description+s.O[i].Brief+"категория:"+s.O[i].Category+" : "+s.O[i].MainCategory)
}

func (s *Goods) Len () int {
	return len(s.O)
}
