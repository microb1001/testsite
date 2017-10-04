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

func Pager (PageMin, pageMax, pageCurrent int, UrlPart string) (newP PagerType) {
	if pageCurrent>PageMin {newP.Prev=UrlPart+"?p="+strconv.Itoa(pageCurrent-1)}
	if pageCurrent<pageMax {newP.Next=UrlPart+"?p="+strconv.Itoa(pageCurrent+1)}
	for ii:=MinMax(pageCurrent-2,0, pageMax);ii<=MinMax(pageCurrent+2,0, pageMax);ii++{
		newP.Elem=append(newP.Elem,PagerElemType{ii+1,"",UrlPart+"?p="+strconv.Itoa(ii),ii == pageCurrent} )
	}
	return
}