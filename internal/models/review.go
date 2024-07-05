package models

import (
	"gorm.io/gorm"
	"time"
)

type Review struct {
	gorm.Model
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	CourseID   uint      `gorm:"not null"`
	Rating     int       `gorm:"not null"`
	Comment    string    `gorm:"type:text"`
	ReviewDate time.Time `gorm:"not null"`
	User       *User     `gorm:"foreignKey:UserID"`   // Belongs To User
	Course     *Course   `gorm:"foreignKey:CourseID"` // Belongs To Course
}

func (Review) TableName() string {
	return "reviews"
}
