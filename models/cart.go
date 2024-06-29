package models

type CartId int

type CartItem Product

type Cart struct {
	Id    CartId
	Items []CartItem
}
