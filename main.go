package main

import (
	"fmt"
	"os"

	"github.com/LucasMRC/lb_back/pkg/auth"
	"github.com/LucasMRC/lb_back/pkg/tasks"
	"github.com/LucasMRC/lb_back/pkg/users"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	fmt.Println("🚀 Starting server")
	godotenv.Load()
	gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(CORSMiddleware())

	// Auth routes
	r.POST("/login", auth.Login)
	r.GET("/logout", auth.Logout)
	r.POST("/register", auth.Register)

	// Task routes
	r.POST("/tasks" /* auth.AuthMiddleware, */, tasks.CreateTask)
	r.GET("/tasks" /* auth.AuthMiddleware, */, tasks.GetTasks)
	r.PATCH("/tasks/:taskId" /* auth.AuthMiddleware, */, tasks.UpdateTask)
	r.DELETE("/tasks/:taskId" /* auth.AuthMiddleware, */, tasks.DeleteTask)

	// User routes
	r.GET("/users", users.GetUsers)
	r.GET("/users/:userId", users.GetUser)
	r.PATCH("/users/:userId", users.UpdateUser)

	// Start server
	port := os.Getenv("PORT")
	fmt.Println("🚀 Up & Running at port", port)
	r.Run(":" + port)
}
