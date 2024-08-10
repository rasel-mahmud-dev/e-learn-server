package courseRoute

import (
	courseHandler "e-learn/internal/handlers/course"
	reviewHandler "e-learn/internal/handlers/review"
	"e-learn/internal/middleware"
	"github.com/gin-gonic/gin"
)

func CourseRoute(r *gin.Engine) {

	r.POST("/api/v1/courses", middleware.AuthenticateMiddleware, courseHandler.CreateCourse)
	r.GET("/api/v1/courses", middleware.AuthenticateMiddleware, courseHandler.GetInstructorCourses)
	r.GET("/api/v1/courses2", courseHandler.GetCourses)
	r.POST("/api/v1/courses/add-review", middleware.AuthenticateMiddleware, reviewHandler.CreateCourseReview)
	r.GET("/api/v1/courses/reviews/:courseId", middleware.AuthenticateMiddleware, reviewHandler.GetCourseReviews)
	r.POST("/api/v1/courses/search", middleware.AuthenticateMiddleware, courseHandler.SearchCourse)
	r.POST("/api/v1/courses/search-suggestion", middleware.AuthenticateMiddleware, courseHandler.GetUserSearchSuggestions)

	//r.GET("/api/v1/courses", courseHandler.GetInstructorCourses)
	//r.GET("/users/:id", handlers.GetUserByID)
	//r.PUT("/users/:id", handlers.UpdateUser)
	//r.DELETE("/users/:id", handlers.DeleteUser)
}
