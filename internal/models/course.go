package models

import (
	"gorm.io/gorm"
	"time"
)

type Course struct {
	gorm.Model
	ID          uint      ` json:"iD" gorm:"primaryKey"`
	Title       string    ` json:"title" gorm:"not null"`
	Description string    ` json:"description" gorm:"type:text;not null"`
	AuthorID    uint      ` json:"authorId" gorm:"not null;foreignKey:UserID"`
	PublishDate time.Time ` json:"publishDate" gorm:""`
	CreatedAt   time.Time ` json:"createdAt" json:"CreatedAt" gorm:"not null"`
	UpdatedAt   time.Time ` json:"updatedAt" json:"UpdatedAt" gorm:"not null"`
	Price       float64   ` json:"price" gorm:"not null"`
	Category    string    ` json:"category" gorm:"not null"`
	Users       []*User   ` json:"users" gorm:"many2many:user_courses;"` // Many-to-many relationship with User
	Reviews     []*Review `json:"reviews" gorm:"foreignKey:CourseID"`    // One-to-many relationship with Review
}
