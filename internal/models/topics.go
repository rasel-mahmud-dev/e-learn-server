package models

import (
	"database/sql"
	"e-learn/internal/database"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Topics struct {
	ID        uint           ` json:"id" gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	CreatedAt time.Time      ` json:"createdAt,omitempty" json:"CreatedAt" gorm:"not null"`
	UpdatedAt time.Time      ` json:"updatedAt,omitempty" json:"UpdatedAt" gorm:"not null"`
	Title     string         `json:"title" gorm:"unique;not null"`
	Slug      string         `json:"slug" gorm:"unique;not null"`
}

func GetOne(c *gin.Context, columns []string, scanFunc func(*sql.Row, *Topics) error, where string, values []any) (*Topics, error) {
	query := fmt.Sprintf("SELECT %s FROM categories %s ", strings.Join(columns, ", "), where)

	user := &Topics{}
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
