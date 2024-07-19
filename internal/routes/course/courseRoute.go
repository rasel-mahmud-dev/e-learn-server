package courseRoute

import (
	courseHandler "e-learn/internal/handlers/course"
	"e-learn/internal/middleware"
	"github.com/gin-gonic/gin"
)

func CourseRoute(r *gin.Engine) {

	r.POST("/api/v1/courses", middleware.AuthenticateMiddleware, courseHandler.CreateCourse)
	r.GET("/api/v1/courses", middleware.AuthenticateMiddleware, courseHandler.GetInstructorCourses)
	//r.GET("/api/v1/courses", courseHandler.GetInstructorCourses)

	//r.GET("/users/:id", handlers.GetUserByID)
	//r.PUT("/users/:id", handlers.UpdateUser)
	//r.DELETE("/users/:id", handlers.DeleteUser)
}
