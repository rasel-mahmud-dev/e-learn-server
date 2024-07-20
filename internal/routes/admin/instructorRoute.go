package adminRoute

import (
	"e-learn/internal/handlers"
	"e-learn/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.Engine) {
	r.GET("/api/v1/instructor/all", middleware.AuthenticateMiddleware, handlers.GetInstructorList)
	r.POST("/api/v1/account/unblock/:accountId/:statusId", middleware.AuthenticateMiddleware, handlers.UnlockAccount)

}
