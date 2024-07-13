package handlers

import (
	"e-learn/internal/database"
	"e-learn/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsers(c *gin.Context) {
	rows, err := database.DB.QueryContext(c, "SELECT id, created_at, updated_at, deleted_at, username, email, password_hash, registration_date, last_login, avatar FROM public.users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.RegistrationDate,
			&user.LastLogin,
			&user.Avatar,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create response without the PasswordHash field
	var response []models.User
	for _, user := range users {
		response = append(response, models.User{
			ID:               user.ID,
			CreatedAt:        user.CreatedAt,
			UpdatedAt:        user.UpdatedAt,
			DeletedAt:        user.DeletedAt,
			Username:         user.Username,
			Email:            user.Email,
			RegistrationDate: user.RegistrationDate,
			LastLogin:        user.LastLogin,
			Avatar:           user.Avatar,
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetUsersProfile(c *gin.Context) {
	//
	//profileId := c.Param("profileId")
	//
	var profile models.Profile
	//
	//result := database.DB.Table("profiles").Where("user_id = ?", profileId).First(&profile)
	//
	//if result.Error != nil {
	//	response.ErrorResponse(c, result.Error, nil)
	//	return
	//}
	//
	//profile.User = nil
	//
	c.JSON(http.StatusOK, profile)

}

func CreateUser(c *gin.Context) {
	//var newUser models.User
	//
	//// Bind JSON or form data
	//if err := c.ShouldBindJSON(&newUser); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//newUser.RegistrationDate = time.Now()
	//
	//// Save to database
	//result := database.DB.Create(&newUser)
	//if result.Error != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	//	return
	//}

	models.GetUsersBySelect(database.DB, []string{"id", "email", "username", "avatar"})

	// Return success response
	c.JSON(http.StatusCreated, "hi")
}

func UpdateProfile(c *gin.Context) {
	//var payload models.Profile
	//
	//if err := c.ShouldBindJSON(&payload); err != nil {
	//	response.ErrorResponse(c, err, nil)
	//	return
	//}
	//
	//authUser := utils.GetAuthUser(c)
	//if authUser == nil {
	//	response.ErrorResponse(c, errors.New("unauthorization"), nil)
	//	return
	//}
	//
	//payload.UserId = authUser.UserId
	//
	//var existingProfile models.Profile
	//result := database.DB.Table("profiles").Where("user_id = ?", payload.UserId).First(&existingProfile)
	//
	//if result.Error != nil {
	//	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	//		// Profile does not exist, create a new one
	//		if err := database.DB.Create(&payload).Error; err != nil {
	//			response.ErrorResponse(c, err, nil)
	//			return
	//		}
	//	} else {
	//		response.ErrorResponse(c, result.Error, nil)
	//		return
	//	}
	//} else {
	//	payload.ID = existingProfile.ID // Ensure we're updating the correct record
	//	if err := database.DB.Save(&payload).Error; err != nil {
	//		response.ErrorResponse(c, err, nil)
	//		return
	//	}
	//}
	//
	//c.JSON(http.StatusCreated, payload)
}

func UpdateProfilePhoto(c *gin.Context) {
	//
	//authUser := utils.GetAuthUser(c)
	//if authUser == nil {
	//	response.ErrorResponse(c, errors.New("unauthorization"), nil)
	//	return
	//}
	//
	//// Parse multipart form
	//err := c.Request.ParseMultipartForm(10 << 20) // 10 MB max
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//// Retrieve file from form data
	//file, handler, err := c.Request.FormFile("avatar")
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "avatar file missing"})
	//	return
	//}
	//defer file.Close()
	//uploadResult := fileUpload.UploadImage2(file, handler.Filename)
	//if uploadResult == nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//var profilePayload models.User
	//profilePayload.ID = authUser.UserId
	//profilePayload.Avatar = uploadResult.SecureURL
	//
	//if err := database.DB.Save(&profilePayload).Error; err != nil {
	//	response.ErrorResponse(c, err, nil)
	//	return
	//}
	//
	//c.JSON(http.StatusCreated, profilePayload)

}
