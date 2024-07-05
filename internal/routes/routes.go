package routes

import (
	"e-learn/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Custom CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5173"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/users", handlers.GetUsers)
	r.POST("/users", handlers.CreateUser)

	r.POST("/sub-categories", handlers.CreateSubCategories)
	r.GET("/sub-categories", handlers.GetSubCategories)

	r.POST("/categories", handlers.CreateCategories)
	r.GET("/categories", handlers.GetCategories)

	r.POST("/topics", handlers.CreateTopics)
	r.GET("/topics", handlers.GetTopics)

	r.POST("/course", handlers.CreateCourse)
	r.GET("/course", handlers.GetCourses)

	//r.GET("/users/:id", handlers.GetUserByID)
	//r.PUT("/users/:id", handlers.UpdateUser)
	//r.DELETE("/users/:id", handlers.DeleteUser)

	return r
}
