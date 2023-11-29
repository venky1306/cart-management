package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/venky1306/cart-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("can't find the product.")
	ErrCantDecodeProducts = errors.New("can't decode the products.")
	ErrUserIdNotValid     = errors.New("user id not valid.")
	ErrCantUpdateUser     = errors.New("can't update user.")
	ErrCantRemoveCartItem = errors.New("can't remove cart item.")
	ErrCantGetItem        = errors.New("can't get item.")
	ErrCantBuyCartItem    = errors.New("can't buy cart item.")
)

func AddProductToCart(ctx context.Context, productCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	var user models.User
	var product models.Product

	filter := bson.M{"_id": productID}
	err := productCollection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return ErrCantFindProduct
	}

	filter = bson.M{"_id": userID}
	err = userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return err
	}

	update := bson.M{"$push": bson.M{"usercart": product}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantFindProduct
	}

	return nil
}

func RemoveCartItem(ctx context.Context, productCollection, userCollection *mongo.Collection, productID primitive.ObjectID, UserID string) error {
	_, err := userCollection.UpdateOne(ctx, bson.M{"user_id": UserID}, bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}})
	if err != nil {
		log.Panic(err)
		return ErrCantRemoveCartItem
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, userID string, userCollection *mongo.Collection) error {
	var order models.Order
	order.Order_ID = primitive.NewObjectID()
	order.Ordered_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		log.Panic(err)
		return ErrUserIdNotValid
	}
	order.Order_Cart = user.UserCart
	user.UserCart = make([]models.Product, 0)
	order.Ordered_To = user.Address_Details
	var cartTotal int = 0
	for _, product := range order.Order_Cart {
		cartTotal += int(*product.Price)
	}
	order.Price = cartTotal
	user.Order_Status = append(user.Order_Status, order)
	user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.D{{"$set", bson.D{{"usercart", user.UserCart}, {"orders", user.Order_Status}, {"updated_at", user.Updated_At}}}})

	if err != nil {
		log.Panic(err)
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuy(ctx context.Context, productCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	var product models.Product
	err := productCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		return ErrCantFindProduct
	}
	var order models.Order
	order.Order_ID = primitive.NewObjectID()
	order.Ordered_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Order_Cart = append(order.Order_Cart, product)
	order.Price = int(*product.Price)

	var user models.User
	userCollection.FindOne(ctx, bson.M{"userid": userID}).Decode(&user)
	order.Ordered_To = user.Address_Details
	user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Order_Status = append(user.Order_Status, order)
	_, err = userCollection.UpdateOne(ctx, bson.M{"userid": userID}, bson.D{{"$set", bson.D{{"updated_at", user.Updated_At}, {"orders", user.Order_Status}}}})
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}
