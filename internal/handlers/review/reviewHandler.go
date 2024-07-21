package reviewHandler

import (
	"database/sql"
	"e-learn/internal/models/review"
	"e-learn/internal/response"
	"e-learn/internal/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	//authUser := utils.GetAuthUser(c)
	//if authUser == nil {
	//	response.ErrorResponse(c, errors.New("unauthorization"), nil)
	//	return
	//}

	orderByQuery := c.Query("orderBy")
	orderQuery := c.Query("order")
	pageNumberQuery := c.Query("pageNumber")

	courseId := c.Param("courseId")
	if courseId == "" {
		response.ErrorResponse(c, errors.New("Missing course id."), nil)
		return
	}

	orderBy := "created_at"
	if orderByQuery != "" {
		if orderByQuery == "date" {
			orderBy = "created_at"
		} else if orderByQuery == "rating" {
			orderBy = "rate"
		}
	}

	limit := 20
	order := "desc"
	if orderQuery != "" {
		if orderQuery == "1" {
			order = "desc"
		} else {
			orderBy = "asc"
		}
	}

	offset := 0
	if pageNumberQuery != "" {
		atoi, err := strconv.Atoi(pageNumberQuery)
		if err != nil {
			// ignore
		}
		offset = limit * (atoi - 1)
	}

	query := fmt.Sprintf("order by %s %s limit 20 offset %d  ", orderBy, order, offset)

	columns := []string{"title", "summary", "rate", "course_id", "user_id", "created_at"}
	info, err := review.GetAllBySelect(c, columns, func(rows *sql.Rows, r *review.Review) error {
		return rows.Scan(&r.Title, &r.Summary, &r.Rate, &r.CourseID, &r.UserID, &r.CreatedAt)
	}, query, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Something were wrong.")})
		return
	}

	rateInfo := review.Review{}
	var oneStar, twoStar, threeStar, fourStar, fiveStar int64

	if offset == 0 {
		columns2 := []string{
			`count(rate)`,
			`ROUND(SUM(rate) / COUNT(rate)) as rate`,
			`(select count(rate) from reviews where rate = 1) as one_star`,
			`(select count(rate) from reviews where rate = 2) as two_star`,
			`(select count(rate) from reviews where rate = 3) as three_star`,
			`(select count(rate) from reviews where rate = 4) as four_star`,
			`(select count(rate) from reviews where rate = 5) as five_sta`,
		}

		rateInfo_, _ := review.GetOne(c, columns2, func(row *sql.Row, r *review.Review) error {
			return row.Scan(&r.Total, &r.Rate, &oneStar, &twoStar, &threeStar, &fourStar, &fiveStar)
		}, "", nil)
		rateInfo = *rateInfo_

	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"reviews":   info,
			"avgRating": rateInfo.Rate,
			"total":     rateInfo.Total,
			"1":         oneStar,
			"2":         twoStar,
			"3":         threeStar,
			"4":         fourStar,
			"5":         fiveStar,
		},
	})

}
