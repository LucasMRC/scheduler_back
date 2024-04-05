package database

import (
	"fmt"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id    int    `json:"id,omitempty"`
	Alias string `json:"alias"`
	Email string `json:"email"`
	Hash  string `json:"-"`
}

func (u *User) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Hash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Hash = string(hashedPassword)
	u.Alias = html.EscapeString(strings.TrimSpace(u.Alias))
	return nil
}

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

type RecurringType struct {
	Id    int
	Title string
}

type TaskStatus struct {
	Id    int
	Title string
}
