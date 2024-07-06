package models

import (
	"time"
)

type Profile struct {
	ID        uint      ` json:"id" gorm:"primaryKey"`
	CreatedAt time.Time ` json:"createdAt,omitempty" json:"CreatedAt" gorm:"not null"`
	UpdatedAt time.Time ` json:"updatedAt,omitempty" json:"UpdatedAt" gorm:"not null"`

	FirstName string `json:"firstName,omitempty" gorm:""`
	LastName  string `json:"lastName,omitempty" gorm:""`
	Headline  string `json:"headline,omitempty" gorm:""`
	Language  string `json:"language,omitempty" gorm:""`
	Website   string `json:"website,omitempty" gorm:""`
	Twitter   string `json:"twitter,omitempty" gorm:""`
	Facebook  string `json:"facebook,omitempty" gorm:""`
	Youtube   string `json:"youtube,omitempty" gorm:""`
	Github    string `json:"github,omitempty" gorm:""`
	AboutMe   string `json:"aboutMe,omitempty" gorm:""`
	UserId    uint   `json:"userId,omitempty" gorm:"unique;primaryKey;index;not null"`
	User      *User  `json:"user,omitempty" gorm:"foreignKey:UserId;references:ID"`
}
