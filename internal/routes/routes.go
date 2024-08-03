package routes

import (
	"e-learn/internal/handlers"
	"e-learn/internal/middleware"
	adminRoute "e-learn/internal/routes/admin"
	authRoute "e-learn/internal/routes/auth"
	courseRoute "e-learn/internal/routes/course"
	instructorRoute "e-learn/internal/routes/instructor"
	roleRoute "e-learn/internal/routes/role"
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
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5173"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong",
			"ClientIp": c.ClientIP(),
			"RemoteIp": c.RemoteIP(),
		})
	})

	r.GET("/users", handlers.GetUsers)

	r.GET("/users/profile/:profileId", middleware.AuthenticateMiddleware, handlers.GetUsersProfile)

	r.PATCH("/users/update-profile", middleware.AuthenticateMiddleware, handlers.UpdateProfile)
	r.PATCH("/users/update-profile-photo", middleware.AuthenticateMiddleware, handlers.UpdateProfilePhoto)

	r.POST("/sub-categories", handlers.CreateSubCategories)
	r.GET("/sub-categories/one", handlers.GetSubCategory)
	r.GET("/sub-categories", handlers.GetSubCategories)
	//r.PATCH("/sub-categories/:id", handlers.UpdateSubCategory)

	r.POST("/categories", handlers.CreateCategories)
	r.GET("/categories", handlers.GetCategories)

	r.POST("/topics", handlers.CreateTopics)
	r.GET("/topics/one", handlers.GetTopic)
	r.GET("/topics", handlers.GetTopics)
	r.PATCH("/topics/:slug", handlers.UpdateTopic)

	//r.POST("/course", handlers.CreateCourse)
	//r.POST("/course/batch", handlers.CreateCourseBatch)
	//r.GET("/course", handlers.GetCourses)
	//
	//r.GET("/course/detail/:slug", handlers.GetCourseDetail)

	instructorRoute.InitRoute(r)
	adminRoute.InitRoute(r)
	authRoute.AuthRoute(r)
	courseRoute.CourseRoute(r)
	roleRoute.RoleRoutes(r)

	return r
}
