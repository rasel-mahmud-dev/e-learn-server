package authRoute

import (
	authHandler "e-learn/internal/handlers/auth"
	"e-learn/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoute(r *gin.Engine) {

	r.POST("/api/v1/auth/login", authHandler.Login)
	r.GET("/api/v1/auth/verify", middleware.AuthenticateMiddleware, authHandler.VerifyUser)

	r.POST("/api/v1/auth/signup", authHandler.CreateUser)

	//r.GET("/users/:id", handlers.GetUserByID)
	//r.PUT("/users/:id", handlers.UpdateUser)
	//r.DELETE("/users/:id", handlers.DeleteUser)
}
