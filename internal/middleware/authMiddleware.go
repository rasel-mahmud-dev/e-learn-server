package middleware

import (
	"e-learn/internal/response"
	"e-learn/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
)

func AuthenticateMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		response.ErrorResponse(c, errors.New("missing Token"), nil)
		return
	}

	data := utils.ParseToken(tokenString)
	c.Set("auth", data)
	c.Next()
}
