package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/venky1306/cart-management/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.GET("/products", controllers.SearchProduct())
	incomingRoutes.GET("/products/search", controllers.SearchProductByQuery())
}
