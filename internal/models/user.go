package models

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID               *int64     `json:"id,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
	Username         string     `json:"username,omitempty"`
	Email            string     `json:"email,omitempty"`
	PasswordHash     string     `json:"password_hash,omitempty"`
	RegistrationDate time.Time  `json:"registration_date,omitempty"`
	LastLogin        *time.Time `json:"last_login,omitempty"`
	Avatar           *string    `json:"avatar,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func GetLoggedUserInfo(userId uint) *User {
	var user User
	//database.DB.First(&user).Select("id,  fullName, username, email").Where("id = ?", userId)
	return &user
}

func GetUsersBySelect(db *sql.DB, columns []string) ([]User, error) {
	table := "users"
	validTables := map[string]bool{table: true}
	validColumns := map[string]bool{
		"id":                true,
		"created_at":        true,
		"updated_at":        true,
		"deleted_at":        true,
		"username":          true,
		"email":             true,
		"password_hash":     true,
		"registration_date": true,
		"last_login":        true,
		"avatar":            true,
	}

	if !validTables[table] {
		return nil, fmt.Errorf("invalid table name")
	}

	for _, col := range columns {
		if !validColumns[col] {
			return nil, fmt.Errorf("invalid column name: %s", col)
		}
	}

	// Build the query string
	//query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns, ", "), table)
	//
	//// Execute the query
	//rows, err := db.Query(query)
	//if err != nil {
	//	return nil, err
	//}
	//defer rows.Close()
	//
	var users []User
	//for rows.Next() {
	//	// Create a slice of interface{} to hold the column values
	//	values := make([]interface{}, len(columns))
	//	valuePtrs := make([]interface{}, len(columns))
	//	for i := range columns {
	//		valuePtrs[i] = &values[i]
	//	}
	//
	//	// Scan the row into the value pointers
	//	if err := rows.Scan(valuePtrs...); err != nil {
	//		return nil, err
	//	}
	//
	//	// Create a User instance and populate it
	//	user := User{}
	//	for i, col := range columns {
	//		field := reflect.ValueOf(&user).Elem().FieldByNameFunc(func(fieldName string) bool {
	//			return strings.EqualFold(fieldName, col)
	//		})
	//		if field.IsValid() && field.CanSet() {
	//			field.Set(reflect.ValueOf(values[i]))
	//		}
	//	}
	//	users = append(users, user)
	//}
	//
	return users, nil
}
