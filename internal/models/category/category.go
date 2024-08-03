package category

import (
	"database/sql"
	"e-learn/internal/database"
	"e-learn/internal/utils"
	"errors"
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
	Title       string     `json:"title,omitempty"`
	Slug        string     `json:"slug,omitempty"`
	Type        *string    `json:"type,omitempty"` // category / subcategory // topic
	Image       *string    `json:"image,omitempty"`
	Description *string    `json:"description,omitempty"`
}

type CategoryWithCamelCaseJSON struct {
	Category
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`

	SubCategoryIds *[]string `json:"subCategories,omitempty"`
}

func GetAllBySelect(c *gin.Context, columns []string, scanFunc func(*sql.Rows, *CategoryWithCamelCaseJSON) error, where string) ([]CategoryWithCamelCaseJSON, error) {
	table := "categories"

	// Build the query string
	query := fmt.Sprintf("SELECT %s FROM %s %s", strings.Join(columns, ", "), table, where)

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

func BatchInsert(c *gin.Context, categories []CategoryWithCamelCaseJSON, asType string) error {

	for _, category := range categories {
		// Build the query string
		query := "insert into categories(title, slug, image, description, created_at, type) values ($1, $2, $3, $4, $5, $6)"

		result, err := database.DB.ExecContext(
			c,
			query,
			category.Title,
			utils.Slugify(category.Title),
			category.Image,
			category.Description,
			time.Now(),
			asType,
		)

		if err != nil {
			fmt.Println(err)
			fmt.Println("errllsdjflskdfjlsdkf")
		}

		fmt.Println(result.RowsAffected())
	}

	fmt.Println("hi")

	return nil

}

func GetOne(c *gin.Context, columns []string, scanFunc func(*sql.Row, *CategoryWithCamelCaseJSON) error, where string, values []any) (*CategoryWithCamelCaseJSON, error) {
	query := fmt.Sprintf("SELECT %s FROM categories %s ", strings.Join(columns, ", "), where)

	user := &CategoryWithCamelCaseJSON{}
	row := database.DB.QueryRowContext(c, query, values...)

	// Use the scan function to populate the users struct
	err := scanFunc(row, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No users found
		}
		return nil, err
	}

	return user, nil
}
