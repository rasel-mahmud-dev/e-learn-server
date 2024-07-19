package course

import (
	"gorm.io/gorm"
	"time"
)

type Course struct {
	ID        uint           ` json:"id" gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	CreatedAt time.Time      ` json:"createdAt,omitempty" json:"CreatedAt" gorm:"not null"`
	UpdatedAt time.Time      ` json:"updatedAt,omitempty" json:"UpdatedAt" gorm:"not null"`

	Thumbnail   string     ` json:"thumbnail" gorm:""`
	Title       string     ` json:"title" gorm:"not null"`
	Slug        string     ` json:"slug" gorm:"unique;not null"`
	Description string     ` json:"description,omitempty" gorm:"type:text;not null"`
	AuthorID    uint       ` json:"authorId" gorm:"not null;foreignKey:UserID"`
	PublishDate *time.Time `json:"publishDate,omitempty" gorm:""`

	Price float64 ` json:"price" gorm:"not null"`
}