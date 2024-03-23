package users

import (
	"github.com/LucasMRC/lb_back/pkg/notion"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users, err := notion.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"users": users})
}

func GetUser(c *gin.Context) {
	c.JSON(200, gin.H{"message": "GetUser"})
}

func UpdateUser(c *gin.Context) {
	c.JSON(200, gin.H{"message": "UpdateUser"})
}
