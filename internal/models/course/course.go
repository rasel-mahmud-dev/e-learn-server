package course

import (
	"database/sql"
	"e-learn/internal/database"
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

	Price           float64 ` json:"price" `
	CategoryList    *string ` json:"categoryList,omitempty" `
	SubCategoryList *string ` json:"subCategoryList,omitempty" `
	TopicList       *string ` json:"topicList,omitempty" `
	AuthorList      *string ` json:"authorList,omitempty" `
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
