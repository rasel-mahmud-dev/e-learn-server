package handlers

import (
	"database/sql"
	"e-learn/internal/models/category"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCategories(c *gin.Context) {

	categories, err := category.GetAllBySelect(c, []string{"id", "title", "slug", "image", "description", "created_at"}, func(rows *sql.Rows, category *category.CategoryWithCamelCaseJSON) error {
		return rows.Scan(
			&category.ID,
			&category.Title,
			&category.Slug,
			&category.Image,
			&category.Description,
			&category.CreatedAt,
		)
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

//func GetCategories(c *gin.Context) {
//
//	var users []models.Category
//	result := database.DB.Find(&users)
//	if result.Error != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
//		return
//	}
//
//	var response []models.Category
//	for _, user := range users {
//		response = append(response, models.Category{
//			ID:    user.ID,
//			Title: user.Title,
//			Slug:  user.Slug,
//		})
//	}
//
//	c.JSON(http.StatusOK, response)
//
//}
//
//func GetTopics(c *gin.Context) {
//
//	var users []models.Topics
//	result := database.DB.Find(&users)
//	if result.Error != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
//		return
//	}
//
//	var response []models.Topics
//	for _, user := range users {
//		response = append(response, models.Topics{
//			ID:    user.ID,
//			Title: user.Title,
//			Slug:  user.Slug,
//		})
//	}
//
//	c.JSON(http.StatusOK, response)
//
//}
//
//func CreateSubCategories(c *gin.Context) {
//	var titles []string
//
//	// Bind JSON or form data
//	if err := c.ShouldBindJSON(&titles); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	var subcategories []models.SubCategory
//	for i := range titles {
//		title := titles[i]
//		subcategories = append(subcategories, models.SubCategory{
//			Title: title,
//			Slug:  utils.Slugify(title),
//		})
//	}
//
//	// Batch insert subcategories
//	if err := database.DB.Create(&subcategories).Error; err != nil {
//		panic(err)
//	}
//
//	c.JSON(http.StatusCreated, subcategories)
//
//}

func CreateCategories(c *gin.Context) {
	var categoriesBody []category.CategoryWithCamelCaseJSON

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&categoriesBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.BatchInsert(c, categoriesBody)

	c.JSON(http.StatusCreated, gin.H{
		"data": categoriesBody,
	})
}

//func CreateTopics(c *gin.Context) {
//	var titles []string
//
//	// Bind JSON or form data
//	if err := c.ShouldBindJSON(&titles); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	var topics []models.Topics
//	for i := range titles {
//		title := titles[i]
//		topics = append(topics, models.Topics{
//			Title: title,
//			Slug:  utils.Slugify(title),
//		})
//
//	}
//
//	// Batch insert topics
//	if err := database.DB.Create(&topics).Error; err != nil {
//
//		response.ErrorResponse(c, err, map[string]string{
//			"unique": "Item already exist.",
//		})
//		return
//	}
//
//	c.JSON(http.StatusCreated, topics)
//
//}
