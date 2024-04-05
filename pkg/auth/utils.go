package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func validateUserInput(c *gin.Context, input *UserInput) error {
	if err := c.ShouldBindJSON(input); err != nil {
		var errorMessage string
		if strings.Contains(err.Error(), "LoginInput.Username") {
			errorMessage = "Invalid input: username is required"
		} else if strings.Contains(err.Error(), "LoginInput.Password") {
			errorMessage = "Invalid input: password is required"
		} else {
			errorMessage = err.Error()
		}
		return errors.New(errorMessage)
	}
	return nil
}

func verifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func generateToken(input UserInput) (string, error) {
	expirationTime := time.Now().Add(8 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Username: input.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	})

	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
