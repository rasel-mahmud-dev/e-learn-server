package review

import (
	"e-learn/internal/database"
	"github.com/gin-gonic/gin"
)

type Review struct {
	ID        uint64  `json:"id,omitempty"`
	Title     string  `json:"title,omitempty"`
	Summary   *string `json:"summary,omitempty"`
	UserID    string  `json:"userId,omitempty"`
	CourseID  string  `json:"courseId,omitempty"`
	Rate      int8    `json:"rate,omitempty"`
	CreatedAt string  `json:"createdAt,omitempty"`
	DeletedAt *string `json:"deletedAt,omitempty"`
}

func InsertOne(c *gin.Context, payload *Review) error {
	query := `insert into reviews(title, summary, course_id, user_id, rate)
		values($1, $2, $3, $4, $5)`

	_, err := database.DB.ExecContext(c,
		query, payload.Title,
		payload.Summary,
		payload.CourseID,
		payload.UserID,
		payload.Rate,
	)
	if err != nil {
		return err
	}

	return nil
}
