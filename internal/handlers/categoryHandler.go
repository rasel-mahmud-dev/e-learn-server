package handlers

import (
	"e-learn/internal/database"
	"e-learn/internal/models"
	"e-learn/internal/response"
	"e-learn/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSubCategories(c *gin.Context) {

	var users []models.SubCategory
	result := database.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var response []models.SubCategory
	for _, user := range users {
		response = append(response, models.SubCategory{
			ID:    user.ID,
			Title: user.Title,
		})
	}

	c.JSON(http.StatusOK, response)

}

func GetCategories(c *gin.Context) {

	var users []models.Category
	result := database.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var response []models.Category
	for _, user := range users {
		response = append(response, models.Category{
			ID:    user.ID,
			Title: user.Title,
			Slug:  user.Slug,
		})
	}

	c.JSON(http.StatusOK, response)

}

func GetTopics(c *gin.Context) {

	var users []models.Topics
	result := database.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var response []models.Topics
	for _, user := range users {
		response = append(response, models.Topics{
			ID:    user.ID,
			Title: user.Title,
			Slug:  user.Slug,
		})
	}

	c.JSON(http.StatusOK, response)

}

func CreateSubCategories(c *gin.Context) {
	var titles []string

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&titles); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var subcategories []models.SubCategory
	for i := range titles {
		title := titles[i]
		subcategories = append(subcategories, models.SubCategory{
			Title: title,
			Slug:  utils.Slugify(title),
		})
	}

	// Batch insert subcategories
	if err := database.DB.Create(&subcategories).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, subcategories)

}

func CreateCategories(c *gin.Context) {
	var titles []string

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&titles); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var categories []models.Category
	for i := range titles {
		title := titles[i]
		categories = append(categories, models.Category{
			Title: title,
			Slug:  utils.Slugify(title),
		})

	}

	// Batch insert subcategories
	if err := database.DB.Create(&categories).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, categories)

}

func CreateTopics(c *gin.Context) {
	var titles []string

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&titles); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var topics []models.Topics
	for i := range titles {
		title := titles[i]
		topics = append(topics, models.Topics{
			Title: title,
			Slug:  utils.Slugify(title),
		})

	}

	// Batch insert topics
	if err := database.DB.Create(&topics).Error; err != nil {

		response.ErrorResponse(c, err, map[string]string{
			"unique": "Item already exist.",
		})
		return
	}

	c.JSON(http.StatusCreated, topics)

}
