package models

import (
	"e-learn/internal/models/users"
	"gorm.io/gorm"
	"time"
)

type Review struct {
	ID         uint           ` json:"id" gorm:"primaryKey"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	CreatedAt  time.Time      ` json:"createdAt,omitempty" json:"CreatedAt" gorm:"not null"`
	UpdatedAt  time.Time      ` json:"updatedAt,omitempty" json:"UpdatedAt" gorm:"not null"`
	UserID     uint           `gorm:"not null"`
	CourseID   uint           `gorm:"not null"`
	Rating     int            `gorm:"not null"`
	Comment    string         `gorm:"type:text"`
	ReviewDate time.Time      `gorm:"not null"`
	User       *users.User    `gorm:"foreignKey:UserID"`   // Belongs To User
	Course     *Course        `gorm:"foreignKey:CourseID"` // Belongs To Course
}

func (Review) TableName() string {
	return "reviews"
}
