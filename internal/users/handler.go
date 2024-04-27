package users

import (
	"fmt"

	"github.com/LucasMRC/lb_back/internal/database"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users, err := database.GetUsers()
	if err != nil {
		fmt.Println("⚠️ Error while getting users: ", err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"users": users})
}

func GetUser(c *gin.Context) {
	user, err := database.GetUser(c.Param("alias"))
	if err != nil {
		fmt.Println("⚠️ Error while getting user: ", err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": user})
}
