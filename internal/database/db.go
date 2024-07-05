package database

import (
	"e-learn/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(postgres.Open(config.Cfg.DATABASE_URI), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	DB = db
}

func SeedDb() {
	//DB.AutoMigrate(&models.User{})
	//DB.AutoMigrate(&models.Review{})
	//DB.AutoMigrate(&models.Course{})
	//DB.AutoMigrate(&models.Category{})
	//DB.AutoMigrate(&models.SubCategory{})
	//DB.AutoMigrate(&models.Topics{})
}
