package user7

import (
	"../cart7"


)
type User_type struct {
	UID string
	login string
	passhash string
	email string
	cookie string
	cart []cart7.Elem
	shipping_adress1 string
	shipping_adress2 string
	shipping_adress3 string
	payment_info string
}