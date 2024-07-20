package handlers

import (
	"database/sql"
	role "e-learn/internal/constant"
	"e-learn/internal/database"
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
)

func GetUsers(c *gin.Context) {

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

func GetInstructorList(c *gin.Context) {

	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("unauthorization"), nil)
		return
	}

	columns := []string{
		"users.user_id",
		"username",
		"email",
		"users.created_at",
		`(select jsonb_agg(jsonb_build_object(
			'status', ass.status,
			'is_status_active', ass.is_status_active,
			'created_at', ass.created_at,
			'note', ass.note,
			'id', ass.id
			)
			) AS account_status    
			from account_status ass 
			where ass.account_id = users.user_id AND is_status_active = true)
		`,
	}

	join := `join users_roles ur on  ur.user_id = users.user_id 
				join roles r on r.role_id = ur.role_id 
					where r.slug = $1 
`

	users, err := users.GetAllBySelect(c, columns, func(rows *sql.Rows, user *users.User) error {
		return rows.Scan(&user.UserID, &user.Username, &user.Email, &user.CreatedAt, &user.AccountStatus)
	},
		join,
		[]any{role.Instructor},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": users,
	})

}

func UnlockAccount(c *gin.Context) {

	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("unauthorization"), nil)
		return
	}

	accountId := c.Param("accountId")
	statusId := c.Param("statusId")
	query := "update account_status set is_status_active = false where account_status.account_id = $1 AND account_status.id = $2"
	_, err := database.DB.ExecContext(c, query, accountId, statusId)
	if err != nil {
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": "Successfully unlocked account",
	})

}
