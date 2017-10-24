package cart7

import (
	"../good7"
)

type Elem struct {
VendorCode1 string
Quantity int
good7.Elem
}

type Cart [] Elem

type Carts map[uint64]Cart