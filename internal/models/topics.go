package models

import "gorm.io/gorm"

type Topics struct {
	gorm.Model
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title" gorm:"unique;not null"`
	Slug  string `json:"slug" gorm:"unique;not null"`
}
