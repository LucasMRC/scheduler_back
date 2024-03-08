package scheduler

import (
	"fmt"
	"io"
	"net/http"
    "encoding/json"

	"github.com/LucasMRC/lb_back/pkg/auth"
	"github.com/LucasMRC/lb_back/pkg/database"
	"github.com/LucasMRC/lb_back/pkg/tasks"
	"github.com/gin-gonic/gin"
)

type NewTaskInput struct {
    Title       string `json:"title" binding:"required"`
    Description *string `json:"description"`
    DueDate     string `json:"due_date" binding:"required"`
    Recurring   *string `json:"recurring"`
    AssignedTo  string `json:"assigned_to" binding:"required"`
}

func CreateTask(c *gin.Context) {
    var input NewTaskInput

    if err := validateTaskInput(c, &input); err != nil {
        fmt.Println("⚠️ Error while validating input: ", err.Error())
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var description string
    if input.Description == nil {
        description = ""
    } else {
        description = *input.Recurring
    }
    var recurring string
    if input.Recurring == nil {
        recurring = "false"
    } else {
        recurring = *input.Recurring
    }
    
    task := tasks.Task{
        AssignedTo: input.AssignedTo,
        DueDate: input.DueDate,
        Title: input.Title,
        Description: description,
        Recurring: recurring,
        Status: tasks.Status.Pending,
    }

    if err := database.CreateTask(task); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating task"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"status": "Task created"})
}

func GetTasks(c *gin.Context) {
    token := c.Request.Header.Get("Authorization")
    username, err := auth.GetUsernameFromToken(token)
    if err != nil {
        fmt.Println("Error getting the username: ", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching tasks"})
        return
    }

    userTasks, err := database.GetTasks(username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching tasks"})
        return
    }

    response := map[string][]tasks.Task{
        "tasks": userTasks,
    }

    c.JSON(http.StatusOK, response)
}

func UpdateTask(c *gin.Context) {
    token := c.Request.Header.Get("Authorization")
    _, err := auth.GetUsernameFromToken(token)
    if err != nil {
        fmt.Println("Error getting the username: ", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating task"})
        return
    }
    taskId := c.Param("taskId")
    jsonBody, err := io.ReadAll(c.Request.Body)
    if err != nil {
        fmt.Println("Error getting the body from request: ", err.Error())
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
        return
    }
    body := tasks.Task{}
    err = json.Unmarshal([]byte(jsonBody), &body)
    if err != nil {
        fmt.Println("Error parsing json to map: ", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating task"})
        return
    }
    task, err := database.UpdateTask(taskId, body)
    if err != nil {
        fmt.Println("Error updating task: ", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating task"})
        return
    }

    response := map[string]tasks.Task{
        "task": task.(tasks.Task),
    }

    c.JSON(http.StatusOK, response)
}

func DeleteTask(c *gin.Context) {
    token := c.Request.Header.Get("Authorization")
    _, err := auth.GetUsernameFromToken(token)
    if err != nil {
        fmt.Println("Error getting the username: ", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching tasks"})
        return
    }
    taskId := c.Param("taskId")

    if err := database.DeleteTask(taskId); err != nil {
        fmt.Println("Error deleting task: ", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching tasks"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "Task deleted"})
}
