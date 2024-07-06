package models

import (
	"gorm.io/gorm"
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
