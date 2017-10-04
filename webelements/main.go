package webelements

import (
	"strconv"
)

type PagerElemType struct{
	Page int
	Class string
	Url string
	Current bool
}

type PagerType struct{
	Elem [] PagerElemType
	Next string
	Prev string
}

func MinMax(index,min,max int) int{
	if index<min {
		return min
	}
	if index>max {
		return max
	}
	return index
}

func Pager (Page, items_per_page, itemsCnt int, urlPart string) (newP PagerType,i,j int) {
	maxPage :=(itemsCnt-1)/ items_per_page // начинается с нуля
	if Page > 0 {newP.Prev= urlPart +"p="+strconv.Itoa(Page-1)}
	if Page < maxPage {newP.Next= urlPart +"p="+strconv.Itoa(Page+1)}
	for ii:=MinMax(Page-2,0, maxPage);ii<=MinMax(Page+2,0, maxPage);ii++{
		newP.Elem=append(newP.Elem,PagerElemType{ii+1,"", urlPart +"p="+strconv.Itoa(ii),ii == Page} )
	}
	i=MinMax((Page)*items_per_page,0, itemsCnt)
	j=MinMax((Page+1)*items_per_page,0, itemsCnt)
	return
}