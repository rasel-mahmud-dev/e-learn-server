package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var secretKey = []byte("your-secret-key")

type JwtPayload struct {
	Email  string `json:"email"`
	UserId uint   `json:"id"`
}

type MyCustomClaims struct {
	Data JwtPayload
	jwt.RegisteredClaims
}

func CreateToken(payload JwtPayload) (string, error) {
	claims := MyCustomClaims{
		payload,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func VerifyToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return false
	}

	// Check if the token is valid
	if !token.Valid {
		return false
	}

	return true

}

func ParseToken(tokenString string) *JwtPayload {

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok {
		return &claims.Data
	}

	return nil
}
