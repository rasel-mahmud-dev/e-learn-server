package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ErrorResponse(c *gin.Context, err error, errorTextMapping map[string]string) {
	var msg = err.Error()

	if errorTextMapping != nil {
		for key := range errorTextMapping {
			if strings.Contains(msg, key) {
				msg = errorTextMapping[key]
				break
			}
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"errorMessage": msg,
	})
}
