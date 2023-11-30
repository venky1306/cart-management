package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/venky1306/cart-management/controllers"
	"github.com/venky1306/cart-management/database"
	"github.com/venky1306/cart-management/middleware"
	"github.com/venky1306/cart-management/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "products"), database.UserData(database.Client, "users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.POST("/v1/products/addproduct", controllers.ProductViewAdmin())
	router.GET("/v1/addtocart", app.AddToCart())
	router.GET("/v1/removeitem", app.RemoveItem())
	router.GET("/v1/checkout", app.BuyFromCart())
	router.GET("/v1/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
