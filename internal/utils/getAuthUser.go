package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAuthUser(c *gin.Context) *JwtPayload {
	auth, isExists := c.Get("auth")
	if !isExists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errorMessage": "Authentication Failed",
		})
		return nil
	}

	authData, ok := auth.(*JwtPayload)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errorMessage": "Invalid authentication data",
		})
		return nil
	}
	return authData
}
