package review

import (
	"database/sql"
	"e-learn/internal/database"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type Review struct {
	ID        uint64  `json:"id,omitempty"`
	Title     string  `json:"title,omitempty"`
	Summary   *string `json:"summary,omitempty"`
	UserID    string  `json:"userId,omitempty"`
	CourseID  string  `json:"courseId,omitempty"`
	Rate      *int8   `json:"rate,omitempty"`
	Total     *int64  `json:"total,omitempty"`
	CreatedAt string  `json:"createdAt,omitempty"`
	DeletedAt *string `json:"deletedAt,omitempty"`

	// populated field
	Username *string `json:"username,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

func InsertOne(c *gin.Context, payload *Review) (*int64, error) {
	query := `insert into reviews(title, summary, course_id, user_id, rate)
		values($1, $2, $3, $4, $5) returning id`

	result, err := database.DB.ExecContext(c,
		query, payload.Title,
		payload.Summary,
		payload.CourseID,
		payload.UserID,
		payload.Rate,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func GetOne(c *gin.Context, columns []string, scanFunc func(*sql.Row, *Review) error, where string, values []any) (*Review, error) {
	query := fmt.Sprintf("SELECT %s FROM reviews %s ", strings.Join(columns, ", "), where)

	fmt.Println(query)
	user := &Review{}
	row := database.DB.QueryRowContext(c, query, values...)

	err := scanFunc(row, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No users found
		}
		return nil, err
	}

	return user, nil
}

func GetAllBySelect(c *gin.Context, columns []string, scanFunc func(*sql.Rows, *Review) error, where string, values []any) ([]Review, error) {
	table := "reviews"

	// Build the query string
	query := fmt.Sprintf("SELECT %s FROM %s %s", strings.Join(columns, ", "), table, where)
	fmt.Println(query)

	// Execute the query
	rows, err := database.DB.QueryContext(c, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		if err := scanFunc(rows, &review); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}
