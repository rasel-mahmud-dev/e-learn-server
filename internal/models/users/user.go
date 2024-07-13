package users

import (
	"database/sql"
	"e-learn/internal/database"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
	"time"
)

type User struct {
	ID               uint64     `json:"id,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
	Username         string     `json:"username,omitempty"`
	Email            string     `json:"email,omitempty"`
	Password         string     `json:"password,omitempty"`
	RegistrationDate time.Time  `json:"registration_date,omitempty"`
	LastLogin        *time.Time `json:"last_login,omitempty"`
	Avatar           *string    `json:"avatar,omitempty"`
}

type UserWithCamelCaseJSON struct {
	User
	ID        uint64     `json:"userId,omitempty"`
	Username  string     `json:"userName,omitempty"`
	Avatar    *string    `json:"avatar,omitempty"`
	LastLogin *time.Time `json:"lastLogin,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func GetLoggedUserInfo(c *gin.Context, userId uint64) (*User, error) {

	payloadAuthInfo, err := GetUserById(c, []string{"id", "username", "email", "avatar"}, func(row *sql.Row, user *User) error {
		return row.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Avatar,
		)
	}, userId)

	if err != nil {
		return nil, err
	}

	if payloadAuthInfo == nil {
		return nil, err
	}

	return payloadAuthInfo, err
}

func GetUsersBySelect(c *gin.Context, columns []string, scanFunc func(*sql.Rows) error) ([]User, error) {
	table := "users"

	// Build the query string
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns, ", "), table)

	// Execute the query
	rows, err := database.DB.QueryContext(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := scanFunc(rows); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByEmail(c *gin.Context, columns []string, scanFunc func(*sql.Row, *User) error, email string) (*User, error) {
	// Build the query string
	query := fmt.Sprintf("SELECT %s FROM users WHERE email = $1", strings.Join(columns, ", "))

	// Prepare a variable to hold the users data
	user := &User{}

	// Execute the query
	row := database.DB.QueryRowContext(c, query, email)

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

func GetUserById(c *gin.Context, columns []string, scanFunc func(*sql.Row, *User) error, userId uint64) (*User, error) {
	// Build the query string
	query := fmt.Sprintf("SELECT %s FROM users WHERE id = $1", strings.Join(columns, ", "))

	// Prepare a variable to hold the users data
	user := &User{}

	// Execute the query
	row := database.DB.QueryRowContext(c, query, userId)

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

func CreateUser(c *gin.Context, payload *User) (*User, error) {

	// Prepare the SQL statement
	query := `
		INSERT INTO users (username, email, password, registration_date, avatar) 
			VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at, updated_at
`
	var user User
	// Execute the SQL statement
	err := database.DB.QueryRowContext(c, query,
		payload.Username,
		payload.Email,
		payload.Password,
		payload.RegistrationDate,
		payload.Avatar,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateProfilePhoto(c *gin.Context, payload *User) (*User, error) {
	var updates []string
	var values []interface{}
	valueIndex := 1

	v := reflect.ValueOf(payload).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Ptr {
			if !field.IsNil() {
				updates = append(updates, fmt.Sprintf("%s = $%d", fieldType.Name, valueIndex))
				values = append(values, field.Interface())
				valueIndex++
			}
		} else if field.IsValid() && field.IsZero() == false {
			updates = append(updates, fmt.Sprintf("%s = $%d", fieldType.Name, valueIndex))
			values = append(values, field.Interface())
			valueIndex++
		}
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Build the SQL statement
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(updates, ", "), valueIndex)
	values = append(values, payload.ID) // Assume payload.ID holds the users ID

	// Execute the SQL statement
	_, err := database.DB.ExecContext(c, query, values...)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
