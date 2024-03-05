package auth

import (
	"fmt"
	"net/http"

	"github.com/LucasMRC/lb_back/pkg/database"
	"github.com/LucasMRC/lb_back/pkg/users"
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
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    if err := verifyPassword(input.Password, user.Password); err != nil {
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

    if err := database.SaveToken(tokenString); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving token"})
        return
    }

    c.JSON(http.StatusOK, tokenString)
}

func Register(c *gin.Context) {
    input := UserInput{}

    if err := validateUserInput(c, &input); err != nil {
        fmt.Println("⚠️ Error while validating input: ", err.Error())
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if user, _ := database.GetUser(input.Username); user.Username != "" {
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
    user := users.User{
        Username: input.Username,
        Password: hash,
    }

    if err := database.SaveUser(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"status": "User created"})
}
