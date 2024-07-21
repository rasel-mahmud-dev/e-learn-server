package reviewHandler

import (
	"e-learn/internal/models/review"
	"e-learn/internal/response"
	"e-learn/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateCourseReview(c *gin.Context) {

	// check auth
	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("unauthorization"), nil)
		return
	}

	var createReview review.Review
	if err := c.ShouldBindJSON(&createReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createReview.UserID = authUser.UserId

	if err := review.InsertOne(c, &createReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": "Successfully Create Course Review",
	})

}
