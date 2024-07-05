package handlers

import (
	"e-learn/internal/database"
	"e-learn/internal/models"
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

func CreateSubCategories(c *gin.Context) {
	var titles []string

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&titles); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var subcategories []models.SubCategory
	for i := range titles {
		subcategories = append(subcategories, models.SubCategory{Title: titles[i]})
	}

	// Batch insert subcategories
	if err := database.DB.Create(&subcategories).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, subcategories)

}
