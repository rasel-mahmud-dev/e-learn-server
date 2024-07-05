package routes

import (
	"e-learn/internal/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/users", handlers.GetUsers)
	r.POST("/users", handlers.CreateUser)

	r.POST("/course", handlers.CreateCourse)
	r.GET("/course", handlers.GetCourses)

	r.POST("/sub-categories", handlers.CreateSubCategories)
	r.GET("/sub-categories", handlers.GetSubCategories)

	r.POST("/categories", handlers.CreateCategories)
	r.GET("/categories", handlers.GetCategories)

	r.POST("/topics", handlers.CreateTopics)
	r.GET("/topics", handlers.GetTopics)

	//r.GET("/users/:id", handlers.GetUserByID)
	//r.PUT("/users/:id", handlers.UpdateUser)
	//r.DELETE("/users/:id", handlers.DeleteUser)

	return r
}
