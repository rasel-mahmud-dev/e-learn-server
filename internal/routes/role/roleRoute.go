package roleRoute

import (
	roleHandler "e-learn/internal/handlers/role"
	"e-learn/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RoleRoutes(r *gin.Engine) {
	r.POST("/api/v1/roles", middleware.AuthenticateMiddleware, roleHandler.CreateRole)
	r.GET("/api/v1/roles", middleware.AuthenticateMiddleware, roleHandler.GetRoles)
	r.GET("/api/v1/roles/users-roles", middleware.AuthenticateMiddleware, roleHandler.GetUsersRoles)
	r.POST("/api/v1/roles/users-roles/:userId", middleware.AuthenticateMiddleware, roleHandler.SetUserRoles)
}
