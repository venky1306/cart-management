package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name      *string            `json:"first_name" validate:"required,min=2,max=30"`
	Last_Name       *string            `json:"last_name" validate:"required,min=2,max=30"`
	Password        *string            `json:"password" validate:"required,min=6"`
	Email           *string            `json:"email" validate:"email,required"`
	Phone           *string            `json:"phone" validate:"required"`
	User_Type       *string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Token           *string            `json:"token"`
	Refresh_Token   *string            `json:"refresh_token"`
	Created_At      time.Time          `json:"created_at"`
	Updated_At      time.Time          `json:"updated_at"`
	User_ID         string             `json:"user_id"`
	UserCart        []Product          `json:"usercart" bson:"usercart"`
	Address_Details Address            `json:"address" bson:"address"`
	Order_Status    []Order            `json:"orders" bson:"orders"`
}

type Product struct {
	ProductId    primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Price        *uint64            `json:"price"`
	Rating       *uint8             `json:"rating"`
	Image        *string            `json:"image"`
}

type Address struct {
	Address_ID primitive.ObjectID `bson:"_id"`
	Address1   *string            `json:"address1" bson:"address1"`
	Address2   *string            `json:"address2" bson:"address2"`
	City       *string            `json:"city_name" bson:"city_name"`
	State      *string            `json:"state" bson:"state"`
	Pincode    *string            `json:"pin_code" bson:"pin_code"`
}

type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id"`
	Order_Cart     []Product          `json:"usercart" bson:"usercart"`
	Ordered_At     time.Time          `json:"ordered_at" bson:"ordered_at"`
	Ordered_To     Address            `json:"ordered_to" bson:"ordered_to"`
	Price          int                `json:"price" bson:"price"`
	Discount       *int               `json:"discount" bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	PayPal bool `json:"paypal" bson:"paypal"`
	Credit bool `json:"credit" bson:"credit"`
}
