package handlers

import (
	"e-learn/internal/database"
	"e-learn/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetUsers(c *gin.Context) {

	var users []models.User
	result := database.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var response []models.User
	for _, user := range users {
		response = append(response, models.User{
			ID:               user.ID,
			FullName:         user.FullName,
			Username:         user.Username,
			Email:            user.Email,
			PasswordHash:     "", // Omitting passwordHash field
			RegistrationDate: user.RegistrationDate,
			LastLogin:        user.LastLogin,
		})
	}

	c.JSON(http.StatusOK, response)

}

func CreateUser(c *gin.Context) {
	var newUser models.User

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser.RegistrationDate = time.Now()

	// Save to database
	result := database.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, newUser)
}
