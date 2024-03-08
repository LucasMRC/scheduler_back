package main

import (
	"fmt"
	"os"

	"github.com/LucasMRC/lb_back/pkg/auth"
	"github.com/LucasMRC/lb_back/pkg/scheduler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main() {
    godotenv.Load()
    gin.DisableConsoleColor()
    r := gin.Default()

    r.POST("/login", auth.Login)
    r.POST("/register", auth.Register)
    r.POST("/tasks", /* auth.AuthMiddleware, */ scheduler.CreateTask)
    r.GET("/tasks", /* auth.AuthMiddleware, */ scheduler.GetTasks)
    r.PATCH("/tasks/:taskId", /* auth.AuthMiddleware, */ scheduler.UpdateTask)
    r.DELETE("/tasks/:taskId", /* auth.AuthMiddleware, */ scheduler.DeleteTask) 

    port := os.Getenv("PORT") 
    fmt.Println("🚀 Up & Running at port", port)
    r.Run(":" + port)
}
