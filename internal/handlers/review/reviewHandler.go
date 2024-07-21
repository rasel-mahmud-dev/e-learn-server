package reviewHandler

import (
	"database/sql"
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

	insertedId, err := review.InsertOne(c, &createReview)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if insertedId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("review creation fail")})
		return
	}

	columns := []string{"title", "summary", "rate", "course_id", "user_id", "created_at"}
	info, err := review.GetOne(c, columns, func(row *sql.Row, r *review.Review) error {
		return row.Scan(&r.Title, &r.Summary, &r.Rate, &r.CourseID, &r.UserID, &r.CreatedAt)
	}, "", nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Something were wrong.")})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    info,
		"message": "Successfully Create Course Review",
	})
}

func GetCourseReviews(c *gin.Context) {

	// check auth
	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("unauthorization"), nil)
		return
	}

	courseId := c.Param("courseId")
	if courseId == "" {
		response.ErrorResponse(c, errors.New("Missing course id."), nil)
		return
	}

	query := "order by created_at desc limit 10"

	columns := []string{"title", "summary", "rate", "course_id", "user_id", "created_at"}
	info, err := review.GetAllBySelect(c, columns, func(rows *sql.Rows, r *review.Review) error {
		return rows.Scan(&r.Title, &r.Summary, &r.Rate, &r.CourseID, &r.UserID, &r.CreatedAt)
	}, query, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Something were wrong.")})
		return
	}

	rateCol := `ROUND(SUM(rate) / COUNT(rate)) as rate`

	rateInfo, err := review.GetOne(c, []string{rateCol, "count(rate)"}, func(row *sql.Row, r *review.Review) error {
		return row.Scan(&r.Rate, &r.Total)
	}, "", nil)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"reviews":   info,
			"avgRating": rateInfo.Rate,
			"total":     rateInfo.Total,
		},
	})

}
