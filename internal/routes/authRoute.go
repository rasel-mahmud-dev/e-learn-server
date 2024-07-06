package routes

import (
	"e-learn/internal/handlers"
	"e-learn/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoute(r *gin.Engine) {

	r.POST("/api/v1/auth/login", handlers.Login)
	r.GET("/api/v1/auth/verify", middleware.AuthenticateMiddleware, handlers.VerifyUser)

	//r.GET("/users/:id", handlers.GetUserByID)
	//r.PUT("/users/:id", handlers.UpdateUser)
	//r.DELETE("/users/:id", handlers.DeleteUser)
}
