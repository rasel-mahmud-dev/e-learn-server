package course

import (
	"database/sql"
	"e-learn/internal/database"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type Course struct {
	ID        uint64     `json:"id,omitempty"`
	CourseID  string     `json:"courseId,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`

	Thumbnail   string ` json:"thumbnail,omitempty"`
	Title       string ` json:"title"`
	Slug        string ` json:"slug" `
	Description string ` json:"description,omitempty" `

	PublishDate *time.Time `json:"publishDate,omitempty" `

	Price float64 ` json:"price" `

	CategoryListJson *string  ` json:"categoryListJson,omitempty" `
	CategoryList     *[]int64 ` json:"categoryList,omitempty" `

	SubCategoryListJson *string  ` json:"subCategoryListJson,omitempty" `
	SubCategoryList     *[]int64 ` json:"subCategoryList,omitempty"`

	TopicListJson *string  ` json:"topicListJson,omitempty" `
	TopicList     *[]int64 ` json:"topicList,omitempty" `

	AuthorListJson *string  ` json:"authorListJson,omitempty" `
	AuthorList     *[]int64 ` json:"authorList,omitempty" `
}

func GetAllBySelect(c *gin.Context, columns []string, scanFunc func(*sql.Rows, *Course) error, where string, values []any) ([]Course, error) {
	table := "courses"

	// Build the query string
	query := fmt.Sprintf("SELECT %s FROM %s %s", strings.Join(columns, ", "), table, where)
	fmt.Println(query)

	// Execute the query
	rows, err := database.DB.QueryContext(c, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []Course
	for rows.Next() {
		var user Course
		if err := scanFunc(rows, &user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetOne(c *gin.Context, columns []string, scanFunc func(*sql.Row, *Course) error, where string, values []any) (*Course, error) {
	query := fmt.Sprintf("SELECT %s FROM courses %s ", strings.Join(columns, ", "), where)

	fmt.Println(query)
	user := &Course{}
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
