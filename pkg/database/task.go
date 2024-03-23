package database

import "fmt"

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"dueDate"`
	Recurring   string `json:"recurring"`
	CreatedBy   string `json:"createdBy"`
	AssignedTo  string `json:"assignedTo"`
}

type statusValue int

const (
	Pending statusValue = iota
	Done
	Overdue
)

func (s statusValue) String() string {
	switch s {
	case Pending:
		return "pending"
	case Done:
		return "done"
	case Overdue:
		return "overdue"
	default:
		panic("invalid status")
	}
}

type status struct {
	Pending string
	Done    string
	Overdue string
}

var Status = status{
	Pending: fmt.Sprint(statusValue(0)),
	Done:    fmt.Sprint(statusValue(1)),
	Overdue: fmt.Sprint(statusValue(2)),
}
