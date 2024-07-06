package handlers

import (
	"e-learn/internal/database"
	"e-learn/internal/models"
	"e-learn/internal/response"
	"e-learn/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	var newUser models.User

	// Bind JSON or form data
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newUser.Email == "" || newUser.PasswordHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password is nil"})
		return
	}

	var user models.User

	row := database.DB.Table("users").Select("id, email, password_hash, username, full_name").
		Where("email = ?", newUser.Email).
		Row()

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Username, &user.FullName)
	if err != nil {
		response.ErrorResponse(c, err, map[string]string{})
		return
	}

	tokenString, err := utils.CreateToken(utils.JwtPayload{
		Email:  user.Email,
		UserId: user.ID,
	})

	c.JSON(http.StatusOK, gin.H{"user": tokenString})
}

func VerifyUser(c *gin.Context) {
	token := c.GetHeader("Authorization")
	jwtPayload := utils.ParseToken(token)

	c.JSON(http.StatusOK, gin.H{
		"Data": jwtPayload,
	})

}
