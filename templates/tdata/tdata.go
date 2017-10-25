package tdata

import(
	"time"
    "../../mydb/good7"
	"../../mydb/user7"
	"../../webelements/pager"
)

type ListItem struct {
	URL string
	URLtoCart string
	Title string
	Image string
	Brief string
	Description string
}

type List struct{
	Title, Body string
	Links       []ListItem
	Pager       pager.Pager
	Cat         []good7.Category1List_type
	Timer       time.Duration //Timer
	Session     uint64
}


type ProductSpec struct {
	Key, Value string
}

type Product struct {
	Spec1       []ProductSpec
	Title       string
	URL         string
	URLtoCart 	string
	Description string
	Session     uint64
}


type CartElem struct {
	//cart7.Cart
	Title       string
	URL         string
	URLtoCart   string
	Image       string
	Brief       string
	Description string
	VendorCode	string
	Price		int
}

type Cart struct {
	Title      string
	Body       string
	user7.User_type
	UserCart   []CartElem //cart7.Cart
	Session    uint64
	TotalCount int
	TotalPrice int
}
