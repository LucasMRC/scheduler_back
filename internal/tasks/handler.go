package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/LucasMRC/lb_back/internal/database"
	"github.com/gin-gonic/gin"
)

type NewTaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date" binding:"required"`
	Recurring   string `json:"recurring"`
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
	if input.Description == "" {
		description = ""
	} else {
		description = input.Description
	}
	var recurring string
	if input.Recurring == "" {
		recurring = "false"
	} else {
		recurring = input.Recurring
	}

	userLoggedIn := c.MustGet("user").(database.User)

	userAssigned, err := database.GetUser(input.AssignedTo)
	if err != nil {
		fmt.Println("⚠️ Error while getting user assigned: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating task"})
		return
	}

	task := database.Task{
		TaskCore: database.TaskCore{
			Title:       input.Title,
			Description: description,
			DueDate:     input.DueDate,
			Recurring:   recurring,
		},
		AssignedTo: userAssigned.Id,
		Status:     0,
		CreatedBy:  userLoggedIn.Id,
	}

	if err := database.CreateTask(task); err != nil {
		fmt.Println("⚠️ Error while creating task: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "Task created"})
}

func GetTasks(c *gin.Context) {
	username := c.MustGet("user").(database.User).Alias
	userTasks, err := database.GetTasks(username)
	if err != nil {
		fmt.Println("Error getting tasks: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching tasks"})
		return
	}

	response := map[string][]database.TaskDTO{
		"tasks": userTasks,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
}

func UpdateTask(c *gin.Context) {
	taskId := c.Param("taskId")

	jsonBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("Error getting the body from request: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	body := database.TaskDTO{}
	err = json.Unmarshal([]byte(jsonBody), &body)
	if err != nil {
		fmt.Println("Error parsing json to map: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating task"})
		return
	}
	id, err := strconv.Atoi(taskId)
	if err != nil {
		fmt.Println("Error converting task id to int: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating task"})
		return
	}
	task, err := database.UpdateTask(id, body)
	if err != nil {
		fmt.Println("Error updating task: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	taskId := c.Param("taskId")
	id, err := strconv.Atoi(taskId)
	if err != nil {
		fmt.Println("Error converting task id to int: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting task"})
		return
	}

	if err := database.DeleteTask(id); err != nil {
		fmt.Println("Error deleting task: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Task deleted"})
}
