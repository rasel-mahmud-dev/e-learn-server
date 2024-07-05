package models

import "gorm.io/gorm"

type SubCategory struct {
	gorm.Model
	ID         uint   `json:"id" gorm:"primaryKey"`
	CategoryID uint   `json:"categoryId" gorm:""`
	Title      string `json:"title" gorm:"unique;not null"`
	Slug       string `json:"slug" gorm:"unique;not null"`
}

func (Course) TableName() string {
	return "sub_categories"
}
