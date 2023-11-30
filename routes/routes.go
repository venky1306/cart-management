package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/venky1306/cart-management/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/v1/users/signup", controllers.SignUp())
	incomingRoutes.POST("/v1/users/login", controllers.Login())
	incomingRoutes.GET("/v1/products", controllers.SearchProduct())
	incomingRoutes.GET("/v1/products/search", controllers.SearchProductByQuery())
}
