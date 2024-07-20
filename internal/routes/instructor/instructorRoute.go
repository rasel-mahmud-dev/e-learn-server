package instructorRoute

import (
	courseHandler "e-learn/internal/handlers/course"
	"e-learn/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.Engine) {
	r.POST("/api/v1/instructor/courses", middleware.AuthenticateMiddleware, courseHandler.CreateCourse)
	r.GET("/api/v1/instructor/courses", middleware.AuthenticateMiddleware, courseHandler.GetInstructorCourses)

}
