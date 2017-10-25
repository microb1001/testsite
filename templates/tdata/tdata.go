package tdata

import(
	"time"
    "../../mydb/good7"
	"../../mydb/user7"
	"../../mydb/cart7"
	"../../webelements/pager"
)



type LinkType struct {
	good7.Elem

	URL, URLtoCart, Title, Image string
	Separator3 int
}
type List struct{
	Title, Body string
	Links       []LinkType
	Pager       pager.Pager
	Cat         []good7.Category1List_type
	Timer       time.Duration //Timer
	Session     uint64
}

type ProductSpec struct {
	Key, Value string
}

type Product struct {
	good7.Elem
	Spec1   []ProductSpec
	Title   string
	URL     string
	Session uint64
}

type LinkType1 struct {
	cart7.Cart
	URL, URLtoCart, Title, Image string
}
type Cart struct {
	Title, Body string
	user7.User_type
	UserCart    cart7.Cart
	Session     uint64
	TotalCount int
	TotalPrice int
}
