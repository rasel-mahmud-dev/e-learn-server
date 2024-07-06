package middleware

import (
	"e-learn/internal/response"
	"e-learn/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthenticateMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	parts := strings.Split(tokenString, "Bearer ")
	if len(parts) != 2 {
		response.ErrorResponse(c, errors.New("missing Token"), nil)
		c.Abort()
		return
	}
	tokenString = parts[1]

	if tokenString == "" {
		response.ErrorResponse(c, errors.New("missing Token"), nil)
		c.Abort()
		return
	}

	data := utils.ParseToken(tokenString)
	if data == nil {
		response.ErrorResponse(c, errors.New("unauthorization"), nil)
		c.Abort()
		return
	}

	c.Set("auth", data)
	c.Next()
}
