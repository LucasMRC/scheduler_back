package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/LucasMRC/lb_back/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var sampleSecretKey = []byte("SecretYouShouldHide")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type UserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email,omitempty"`
}

func Login(c *gin.Context) {
	var input UserInput

	if err := validateUserInput(c, &input); err != nil {
		fmt.Println("⚠️ Error while validating input: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := database.GetUser(input.Username)
	if err != nil {
		fmt.Println("⚠️ Error while getting user: ", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err := verifyPassword(input.Password, user.Hash); err != nil {
		fmt.Println("⚠️ Error while verifying password: ", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenString, err := generateToken(input)
	if err != nil {
		fmt.Println("⚠️ Error while generating token: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating token"})
		return
	}

	if err := database.SaveToken(tokenString, user); err != nil {
		fmt.Println("⚠️ Error while saving token: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving token"})
		return
	}

	c.Writer.Header().Set("Access-Control-Expose-Headers", "Authorization")
	c.Writer.Header().Set("Authorization", "Bearer "+tokenString)

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Register(c *gin.Context) {
	input := UserInput{}

	if err := validateUserInput(c, &input); err != nil {
		fmt.Println("⚠️ Error while validating input: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user, _ := database.GetUser(input.Username); user.Alias != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already registered"})
		return
	}

	// Validate email here...

	hash, err := hashPassword(input.Password)
	if err != nil {
		fmt.Println("⚠️ Error while hashing password: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing password"})
		return
	}
	user := database.User{
		Alias: input.Username,
		Hash:  hash,
		Email: input.Email,
	}

	if err := database.SaveUser(user); err != nil {
		fmt.Println("⚠️ Error while saving user: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "User created"})
}

func GetUsernameFromToken(tokenHeader string) (string, error) {
	tokenString := strings.Replace(tokenHeader, "Bearer ", "", 1)
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if ok := token.Method == jwt.SigningMethodHS256; !ok {
			fmt.Println("unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}
		return sampleSecretKey, nil
	})
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}

func Logout(c *gin.Context) {
	tokenHeader := c.GetHeader("Authorization")
	tokenString := strings.Replace(tokenHeader, "Bearer ", "", 1)
	if err := database.DeleteToken(tokenString); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GetSession(c *gin.Context) {
	tokenHeader := c.GetHeader("Authorization")
	tokenString := strings.Replace(tokenHeader, "Bearer ", "", 1)
	user, err := database.GetUserFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
