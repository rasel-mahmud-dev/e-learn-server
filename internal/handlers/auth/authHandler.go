package authHandler

import (
	"database/sql"
	"e-learn/internal/models/users"
	"e-learn/internal/response"
	"e-learn/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Login(c *gin.Context) {
	var newUser users.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newUser.Email == "" || newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password is nil"})
		return
	}

	user, err := users.GetUserByEmail(c, []string{"id", "email", "password", "username", "avatar"}, func(row *sql.Row, user *users.User) error {
		return row.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Username,
			&user.Avatar,
		)
	}, newUser.Email)

	if err != nil {
		response.ErrorResponse(c, err, nil)
		return
	}

	if user == nil {
		response.ErrorResponse(c, err, map[string]string{
			"no rows": "User not registered.",
		})
	}

	if user.Password != newUser.Password {
		response.ErrorResponse(c, errors.New("Wrong Password"), nil)
		return
	}

	tokenString, err := utils.CreateToken(utils.JwtPayload{
		Email:  user.Email,
		UserId: user.ID,
	})

	payloadAuthInfo, err := users.GetLoggedUserInfo(c, user.ID)
	if payloadAuthInfo == nil {
		response.ErrorResponse(c, err, nil)
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "auth": payloadAuthInfo})
}

func VerifyUser(c *gin.Context) {
	authUser := utils.GetAuthUser(c)
	if authUser == nil {
		response.ErrorResponse(c, errors.New("Unauthorization"), nil)
		return
	}

	payloadAuthInfo, err := users.GetLoggedUserInfo(c, authUser.UserId)
	if payloadAuthInfo == nil {
		response.ErrorResponse(c, err, nil)
	}

	c.JSON(http.StatusOK, gin.H{
		"auth": payloadAuthInfo,
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
	newUser.UserID = utils.GenUUID()

	result, error := users.CreateUser(c, &newUser)
	if error != nil {
		response.ErrorResponse(c, error, map[string]string{
			"uni_users_email": "User already registered.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": result})

}
