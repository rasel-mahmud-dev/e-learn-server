package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID               uint       `json:"id" gorm:"primaryKey"`
	FullName         string     `json:"fullName" gorm:"not null"`
	Username         string     `json:"username" gorm:"unique;not null"`
	Email            string     `json:"email" gorm:"unique;not null"`
	PasswordHash     string     `json:"password,omitempty" gorm:"not null"`
	RegistrationDate time.Time  `json:"registrationDate" gorm:"not null"`
	LastLogin        *time.Time `json:"lastLogin" gorm:"default:null"`
	Courses          []*Course  `json:"courses,omitempty" gorm:"many2many:user_courses;"` // Assuming Courses model is defined
	Reviews          []*Review  `json:"reviews,omitempty" gorm:"foreignKey:UserID"`       // Assuming Reviews model is defined
}

func (User) TableName() string {
	return "users"
}
