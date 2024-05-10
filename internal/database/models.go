package database

import (
	"database/sql/driver"
	"errors"
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

type TaskCore struct {
	Id          int        `json:"id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     string     `json:"dueDate"`
	Recurring   string     `json:"recurring"`
	DoneDate    nullString `json:"doneDate"`
}

type TaskDTO struct {
	TaskCore
	Status     string `json:"status"`
	CreatedBy  string `json:"createdBy"`
	AssignedTo string `json:"assignedTo"`
}

type Task struct {
	TaskCore
	Status     int
	CreatedBy  int
	AssignedTo int
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

// Null string to handle the done_date column
type nullString string

func (s *nullString) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	strVal, ok := value.(string)
	if !ok {
		return errors.New("Column is not a string")
	}
	*s = nullString(strVal)
	return nil
}

func (s nullString) Value() (driver.Value, error) {
	if len(s) == 0 { // if nil or empty string
		return nil, nil
	}
	return string(s), nil
}
