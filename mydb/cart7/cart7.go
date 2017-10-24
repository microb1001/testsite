package cart7

import (
	"../good7"
)

type Elem struct {
VendorCode string
Quantity int
good7.Elem
}

type Cart [] Elem