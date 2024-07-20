package role

import (
	"database/sql"
	"e-learn/internal/database"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type Role struct {
	ID          uint64     `json:"id,omitempty"`
	RoleId      string     `json:"roleId,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Status      *string    `json:"status,omitempty"` // active // inactive
	Description *string    `json:"description,omitempty"`
}

type UserRole struct {
	Email  string `json:"email,omitempty"`
	UserId uint64 `json:"userId,omitempty"`
	Roles  string `json:"roles,omitempty"`
}

func GetAllBySelect(c *gin.Context, columns []string, scanFunc func(*sql.Rows, *Role) error, where string) ([]Role, error) {
	table := "roles"

	// Build the query string
	query := fmt.Sprintf("SELECT %s FROM %s %s", strings.Join(columns, ", "), table, where)

	// Execute the query
	rows, err := database.DB.QueryContext(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []Role
	for rows.Next() {
		var user Role
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

func GetUsersRolesMapping(c *gin.Context) ([]UserRole, error) {
	query := `select u.email, u.id, json_agg(DISTINCT  r.role_id) AS "roles"
		from users u
				 left join public.users_roles ur
					  on ur.user_id = u.id
				 left join roles r on ur.role_id = r.role_id
		group by u.email, u.id`

	rows, err := database.DB.QueryContext(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersRoles []UserRole
	for rows.Next() {
		var userRole UserRole
		err := rows.Scan(&userRole.Email, &userRole.UserId, &userRole.Roles)
		if err != nil {
			return nil, err
		}
		usersRoles = append(usersRoles, userRole)
	}

	return usersRoles, nil
}

//func BatchInsert(c *gin.Context, categories []Role, asType string) error {
//
//	for _, category := range categories {
//		// Build the query string
//		query := "insert into categories(title, slug, image, description, created_at, type) values ($1, $2, $3, $4, $5, $6)"
//
//		result, err := database.DB.ExecContext(
//			c,
//			query,
//			category.Title,
//			utils.Slugify(category.Title),
//			category.Image,
//			category.Description,
//			time.Now(),
//			asType,
//		)
//
//		if err != nil {
//			fmt.Println(err)
//			fmt.Println("errllsdjflskdfjlsdkf")
//		}
//
//		fmt.Println(result.RowsAffected())
//	}
//
//	fmt.Println("hi")
//
//	return nil
//
//}

func GetOne(c *gin.Context, columns []string, scanFunc func(*sql.Row, *Role) error, where string, values []any) (*Role, error) {
	query := fmt.Sprintf("SELECT %s FROM categories %s ", strings.Join(columns, ", "), where)

	user := &Role{}
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
