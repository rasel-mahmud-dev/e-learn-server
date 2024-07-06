package models

import (
	"e-learn/internal/database"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           ` json:"id" gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	CreatedAt time.Time      ` json:"createdAt,omitempty" json:"CreatedAt" gorm:"not null"`
	UpdatedAt time.Time      ` json:"updatedAt,omitempty" json:"UpdatedAt" gorm:"not null"`

	FullName         string     `json:"fullName,omitempty" gorm:"not null"`
	Username         string     `json:"username,omitempty" gorm:"unique;not null"`
	Email            string     `json:"email,omitempty" gorm:"unique;not null"`
	PasswordHash     string     `json:"password,omitempty" gorm:"not null"`
	RegistrationDate time.Time  `json:"registrationDate,omitempty" gorm:"not null"`
	LastLogin        *time.Time `json:"lastLogin,omitempty" gorm:"default:null"`
	Courses          []*Course  `json:"courses,omitempty" gorm:"many2many:user_courses;"` // Assuming Courses model is defined
	Reviews          []*Review  `json:"reviews,omitempty" gorm:"foreignKey:UserID"`       // Assuming Reviews model is defined
}

func (User) TableName() string {
	return "users"
}

func GetLoggedUserInfo(userId uint) *User {
	var user User
	database.DB.First(&user).Select("id,  fullName, username, email").Where("id = ?", userId)
	return &user
}
