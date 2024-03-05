package scheduler

import (
	"fmt"
	"net/http"

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
