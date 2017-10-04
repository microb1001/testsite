package webelements

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

func Pager (min, max, curr int, UrlPart string) (newP PagerType) {
	
	return
}