package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("can't find the product.")
	ErrCantDecodeProducts = errors.New("can't decode the products.")
	ErrUserIdNotValid     = errors.New("user id not valid.")
	ErrCantUpdateUser     = errors.New("can't update user.")
	ErrCantRemoveCartItem = errors.New("can't remove cart item.")
	ErrCantGetItem        = errors.New("can't get item.")
	ErrCantBuyCartItem    = errors.New("can't buy cart item.")
)

func AddProductToCart() {

}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuy() {

}
