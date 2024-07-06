package handlers

import (
	"e-learn/internal/database"
	"e-learn/internal/fileUpload"
	"e-learn/internal/models"
	"e-learn/internal/response"
	"e-learn/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func GetUsersProfile(c *gin.Context) {

	profileId := c.Param("profileId")

	var profile models.Profile

	result := database.DB.Table("profiles").Where("user_id = ?", profileId).First(&profile)

	if result.Error != nil {
		response.ErrorResponse(c, result.Error, nil)
		return
	}

	profile.User = nil

	c.JSON(http.StatusOK, profile)

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

func UpdateProfile(c *gin.Context) {
	var payload models.Profile

	if err := c.ShouldBindJSON(&payload); err != nil {
		response.ErrorResponse(c, err, nil)
		return
	}

	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("unauthorization"), nil)
		return
	}

	payload.UserId = authUser.UserId

	var existingProfile models.Profile
	result := database.DB.Table("profiles").Where("user_id = ?", payload.UserId).First(&existingProfile)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Profile does not exist, create a new one
			if err := database.DB.Create(&payload).Error; err != nil {
				response.ErrorResponse(c, err, nil)
				return
			}
		} else {
			response.ErrorResponse(c, result.Error, nil)
			return
		}
	} else {
		payload.ID = existingProfile.ID // Ensure we're updating the correct record
		if err := database.DB.Save(&payload).Error; err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
	}

	c.JSON(http.StatusCreated, payload)
}

func UpdateProfilePhoto(c *gin.Context) {

	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("unauthorization"), nil)
		return
	}

	// Parse multipart form
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve file from form data
	file, handler, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "avatar file missing"})
		return
	}
	defer file.Close()
	uploadResult := fileUpload.UploadImage2(file, handler.Filename)
	if uploadResult == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var profilePayload models.User
	profilePayload.ID = authUser.UserId
	profilePayload.Avatar = uploadResult.SecureURL

	if err := database.DB.Save(&profilePayload).Error; err != nil {
		response.ErrorResponse(c, err, nil)
		return
	}

	c.JSON(http.StatusCreated, profilePayload)

}
