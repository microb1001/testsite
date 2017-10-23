package pager

import(
	"strconv"
	"../servo"
)

type Elem struct{
	Page int
	Class string
	Url string
	Current bool
}

type Pager struct{
	Elem [] Elem
	Next string
	Prev string
}
//
// Возвращает структуру для создания страниц в template
// i,j какие элементы списка сейчас будут отображаться
//
func Set(Page, items_per_page, itemsCnt int, urlPart string) (newP Pager, i, j int) {
	const PAGERWIDTH  = 2
	maxPage := (itemsCnt - 1) / items_per_page // начинается с нуля
	if Page > 0 {
		newP.Prev = urlPart + "p=" + strconv.Itoa(Page-1)
	}
	if Page < maxPage {
		newP.Next = urlPart + "p=" + strconv.Itoa(Page+1)
	}
	for ii := servo.MinMax(Page-PAGERWIDTH, 0, maxPage); ii <= servo.MinMax(Page+PAGERWIDTH, 0, maxPage); ii++ {
		newP.Elem = append(newP.Elem, Elem{ii + 1, "", urlPart + "p=" + strconv.Itoa(ii), ii == Page})
	}
	i = servo.MinMax((Page)*items_per_page, 0, itemsCnt)
	j = servo.MinMax((Page+1)*items_per_page, 0, itemsCnt)
	return
}