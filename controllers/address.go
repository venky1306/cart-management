package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/venky1306/cart-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func EditAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userid := c.Query("userID")
		if userid != "" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "userID missing"})
		}

		var addr models.Address
		if err := c.BindJSON(&addr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		validateErr := validate.Struct(addr)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		var user models.User

		err := UserCollection.FindOne(ctx, bson.M{"_id": userid}).Decode(&user)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		addr.Address_ID = primitive.NewObjectID()
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		result, err := UserCollection.UpdateByID(ctx, bson.D{{Key: "_id", Value: userid}}, bson.D{{Key: "$set", Value: bson.D{{Key: "address", Value: addr}, {Key: "updtaed_at", Value: user.Updated_At}}}})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to add address"})
			return
		}
		c.JSON(http.StatusOK, result)
	}

}
