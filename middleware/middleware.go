package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/venky1306/cart-management/tokens"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("token")
		if clientToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization token missng."})
			return
		}

		var claims *tokens.Claims
		claims, msg := tokens.ValidateToken(clientToken)
		if msg != "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}
		c.Set("email", claims.Email)
		c.Set("firstName", claims.First_Name)
		c.Set("lastName", claims.Last_Name)
		c.Set("userType", claims.User_Type)
		c.Set("userId", claims.User_ID)
		// fmt.Println("********** Abort does not stop this handler function. **********")
		c.Next()
	}
}
