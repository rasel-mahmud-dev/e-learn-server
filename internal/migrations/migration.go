package migrations

import (
	"gorm.io/gorm"
)

func SeedDb(DB *gorm.DB) {
	//DB.AutoMigrate(&models.Course{})
	//DB.Delete(&models.User{})
	//DB.AutoMigrate(&models.User{})
	//DB.AutoMigrate(&models.Review{})
	//DB.AutoMigrate(&models.Course{})
	//DB.AutoMigrate(&models.Category{})
	//DB.AutoMigrate(&models.SubCategory{})
	//DB.AutoMigrate(&models.Topics{})
	//DB.AutoMigrate(&models.Course{})
	//DB.AutoMigrate(&models.Profile{})

	// Add name field
	//DB.Migrator().DropColumn(&models.Course{}, "PublishDate")
	//DB.Migrator().AddColumn(&models.Course{}, "PublishDate")

	//DB.AutoMigrate(&models.Profile{})
}
