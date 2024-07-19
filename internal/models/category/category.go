package category

import (
	"database/sql"
	"e-learn/internal/database"
	"e-learn/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// Category represents the categories table
type Category struct {
	ID          uint64     `json:"id,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Image       *string    `json:"image,omitempty"`
	Description *string    `json:"description,omitempty"`
}

type CategoryWithCamelCaseJSON struct {
	Category
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

func GetAllBySelect(c *gin.Context, columns []string, scanFunc func(*sql.Rows, *CategoryWithCamelCaseJSON) error) ([]CategoryWithCamelCaseJSON, error) {
	table := "categories"

	// Build the query string
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns, ", "), table)

	// Execute the query
	rows, err := database.DB.QueryContext(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []CategoryWithCamelCaseJSON
	for rows.Next() {
		var user CategoryWithCamelCaseJSON
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

func BatchInsert(c *gin.Context, categories []CategoryWithCamelCaseJSON) error {

	for _, category := range categories {
		// Build the query string
		query := "insert into categories(title, slug, image, description, created_at) values ($1, $2, $3, $4, $5)"

		result, err := database.DB.ExecContext(
			c,
			query,
			category.Title,
			utils.Slugify(category.Title),
			category.Image,
			category.Description,
			time.Now(),
		)

		if err != nil {
			fmt.Println("errllsdjflskdfjlsdkf")
		}

		fmt.Println(result.RowsAffected())
	}

	fmt.Println("hi")

	return nil

}
