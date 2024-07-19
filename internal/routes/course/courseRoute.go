package courseRoute

import (
	courseHandler "e-learn/internal/handlers/course"
	"github.com/gin-gonic/gin"
)

func CourseRoute(r *gin.Engine) {

	r.POST("/api/v1/courses", courseHandler.CreateCourse)

	//r.GET("/users/:id", handlers.GetUserByID)
	//r.PUT("/users/:id", handlers.UpdateUser)
	//r.DELETE("/users/:id", handlers.DeleteUser)
}
