package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/venky1306/cart-management/authorization"
	"github.com/venky1306/cart-management/database"
	"github.com/venky1306/cart-management/models"
	"github.com/venky1306/cart-management/tokens"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()
var UserCollection *mongo.Collection = database.UserData(database.Client, "users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "products")

func HashPassowd(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Panic(err)
	}
	return string(hash)
}

func VerifyPassword(userPassword, givenPassword string) bool {
	if bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword)) != nil {
		return true
	}
	return false
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		// fmt.Println("binding works")

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed"})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// fmt.Println("email check passed")
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "phone already exists"})
			return
		}

		password := HashPassowd(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		token, refreshToken, err := tokens.GenerateAllTokens(*user.Email, *user.First_Name, *user.Last_Name, *user.User_Type, user.User_ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating tokens"})
			return
		}

		user.Token = &token
		user.Refresh_Token = &refreshToken
		user.UserCart = make([]models.Product, 0)
		user.Order_Status = make([]models.Order, 0)

		result, err := UserCollection.InsertOne(ctx, user)
		if err != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusCreated, result)
		return

	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		// log.Println(*foundUser.Password, *user.Password)
		isValid := VerifyPassword(*foundUser.Password, *user.Password)
		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
			return
		}

		token, refreshToken, err := tokens.GenerateAllTokens(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, *foundUser.User_Type, foundUser.User_ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating tokens"})
			return
		}

		err = tokens.UpdateTokens(token, refreshToken, foundUser.User_ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		err = UserCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_ID}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

func ProductViewAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// log.Println("middleware works")
		if err := authorization.IsAdmin(c); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var product models.Product
		err := c.BindJSON(&product)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to bind the product info."})
			return
		}

		product.ProductId = primitive.NewObjectID()
		result, err := ProductCollection.InsertOne(ctx, product)
		if err != nil {
			log.Panic(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "product not added."})
			return
		}
		c.IndentedJSON(200, result)
		return
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var products []models.Product

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, "something went wrong. Please try after some time.")
			return
		}

		if err = cursor.All(ctx, &products); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err = cursor.Err(); err != nil {
			log.Println(err)
			c.JSON(400, "invalid")
			return
		}

		c.JSON(200, products)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		queryParam := c.Query("productid")
		if queryParam == "" {
			log.Println("product id is empty")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "productid param missing"})
			return
		}

		var products []models.Product
		cursor, err := ProductCollection.Find(ctx, bson.M{"name": bson.M{"$regex": queryParam}})
		if err != nil {
			c.JSON(http.StatusNotFound, "cannot find requested resource")
			return
		}

		if err = cursor.All(ctx, &products); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)
		if err = cursor.Err(); err != nil {
			log.Println(err)
			c.JSON(400, "invalid")
			return
		}
		c.JSON(200, products)
		return
	}
}
