package authorization

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func IsAdmin(c *gin.Context) error {
	usertype, _ := c.Get("userType")
	if usertype != "ADMIN" {
		return errors.New("Unauthorized to access this resource")
	}
	return nil
}

func AccessUserToUid(c *gin.Context, ToAccessuserId string) error {
	usertype, _ := c.Get("userType")
	if usertype == "ADMIN" || c.GetString("userId") == ToAccessuserId {
		return nil
	}
	return errors.New("Unauthorized to access this resource")
}
