package scheduler

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func validateTaskInput(c *gin.Context, input *NewTaskInput) error {
    if err := c.ShouldBindJSON(input); err != nil {
        var errorMessage string
        if strings.Contains(err.Error(), "NewTaskInput.Title") {
            errorMessage = "Invalid input: title is required"
        } else if strings.Contains(err.Error(), "NewTaskInput.DueDate") {
            errorMessage = "Invalid input: due_date is required"
        } else if strings.Contains(err.Error(), "NewTaskInput.AssignedTo") {
            errorMessage = "Invalid input: assigned_to is required"
        } else {
            errorMessage = err.Error()
        }
        return errors.New(errorMessage)
     }
    return nil
}
