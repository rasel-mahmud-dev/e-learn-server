package handlers

import (
	"database/sql"
	"e-learn/internal/fileUpload"
	"e-learn/internal/models"
	"e-learn/internal/models/users"
	"e-learn/internal/response"
	"e-learn/internal/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func GetUsers(c *gin.Context) {
	fmt.Println("hi")
	users, err := users.GetUsersBySelect(c, []string{"id", "email", "username", "avatar"}, func(rows *sql.Rows, user *users.User) error {
		return rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Avatar,
		)
	})

	if err != nil {
		response.ErrorResponse(c, err, nil)
		return
	}

	c.JSON(http.StatusCreated, users)
}

func GetUsersProfile(c *gin.Context) {

	profileId := c.Param("profileId")

	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("Unauthorization"), nil)
		return
	}

	atoi, err := strconv.Atoi(profileId)
	if err != nil {
		fmt.Println(err)
		return
	}

	myUint64 := uint64(atoi) // Convert int to uint64

	columns := []string{
		"id",
		"created_at",
		"updated_at",
		"first_name",
		"last_name",
		"headline",
		"language",
		"website",
		"twitter",
		"facebook",
		"youtube",
		"github",
		"about_me",
		"user_id",
	}

	payloadAuthInfo, err := models.GetProfileById(c, columns, func(row *sql.Row, profile *models.Profile) error {
		return row.Scan(
			&profile.ID,
			&profile.CreatedAt,
			&profile.UpdatedAt,
			&profile.FirstName,
			&profile.LastName,
			&profile.Headline,
			&profile.Language,
			&profile.Website,
			&profile.Twitter,
			&profile.Facebook,
			&profile.YouTube,
			&profile.Github,
			&profile.AboutMe,
			&profile.UserId,
		)
	}, myUint64)

	if payloadAuthInfo == nil {
		response.ErrorResponse(c, err, nil)
		return
	}

	camelCaseProfile := models.ProfileWithCamelCaseJSON{
		Profile:   *payloadAuthInfo,
		DeletedAt: payloadAuthInfo.DeletedAt,
		CreatedAt: payloadAuthInfo.CreatedAt,
		UpdatedAt: payloadAuthInfo.UpdatedAt,
		FirstName: payloadAuthInfo.FirstName,
		LastName:  payloadAuthInfo.LastName,
		AboutMe:   payloadAuthInfo.AboutMe,
		UserId:    payloadAuthInfo.UserId,
	}

	camelCaseProfile.Profile.AboutMe = nil
	camelCaseProfile.Profile.UpdatedAt = nil
	camelCaseProfile.Profile.UserId = 0
	camelCaseProfile.Profile.FirstName = nil
	camelCaseProfile.Profile.LastName = nil
	camelCaseProfile.Profile.CreatedAt = nil

	c.JSON(http.StatusOK, gin.H{
		"data": camelCaseProfile,
	})
}

func CreateUser(c *gin.Context) {
	var newUser users.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := users.GetUserByEmail(c, []string{"id", "email", "username", "avatar"}, func(row *sql.Row, user *users.User) error {
		return row.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Avatar,
		)
	}, newUser.Email)

	if err != nil {
		response.ErrorResponse(c, err, nil)
		return
	}

	if user != nil {
		response.ErrorResponse(c, errors.New("users is already registered"), nil)
		return
	}

	newUser.RegistrationDate = time.Now()

	result, error := users.CreateUser(c, &newUser)
	if error != nil {
		response.ErrorResponse(c, error, map[string]string{
			"uni_users_email": "User already registered.",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"data": result})

}

func UpdateProfile(c *gin.Context) {
	var payload *models.ProfileWithCamelCaseJSON

	if err := c.ShouldBindJSON(&payload); err != nil {
		response.ErrorResponse(c, err, nil)
		return
	}

	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("unauthorization"), nil)
		return
	}

	payload2 := models.Profile{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Headline:  payload.Headline,
		Language:  payload.Language,
		Website:   payload.Website,
		Twitter:   payload.Twitter,
		Facebook:  payload.Facebook,
		YouTube:   payload.YouTube,
		Github:    payload.Github,
		AboutMe:   payload.AboutMe,
	}

	payload2.UserId = authUser.UserId

	_, err := models.UpdateProfile(c, &payload2)

	if err != nil {
		response.ErrorResponse(c, err, nil)
		return
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

	var profilePayload users.User
	profilePayload.ID = authUser.UserId
	profilePayload.Avatar = utils.StringPtr(uploadResult.SecureURL)

	_, err = users.UpdateProfilePhoto(c, &profilePayload)
	if err != nil {
		return
	}

	//if err := database.DB.Save(&profilePayload).Error; err != nil {
	//	response.ErrorResponse(c, err, nil)
	//	return
	//}

	c.JSON(http.StatusCreated, profilePayload)

}
