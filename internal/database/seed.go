package database

import (
	"database/sql"
	"fmt"
	"time"
)

var tasks = []Task{
	// Tasks for today
	{
		TaskCore: TaskCore{
			Id:          1,
			Title:       "Finish report",
			Description: "Complete the quarterly report for the finance department",
			DueDate:     time.Now().Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     1,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	{
		TaskCore: TaskCore{
			Id:          2,
			Title:       "Meeting with client",
			Description: "Discuss project updates with the client",
			DueDate:     time.Now().Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     3,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	// Tasks for this week
	{
		TaskCore: TaskCore{
			Id:          3,
			Title:       "Prepare presentation",
			Description: "Create slides for the team meeting on Friday",
			DueDate:     time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     1,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	{
		TaskCore: TaskCore{
			Id:          4,
			Title:       "Review project proposal",
			Description: "Provide feedback on the latest project proposal",
			DueDate:     time.Now().AddDate(0, 0, 3).Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     1,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	{
		TaskCore: TaskCore{
			Id:          5,
			Title:       "Monthly review",
			Description: "Conduct monthly performance review with team members",
			DueDate:     time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
			Recurring:   "Monthly",
		},
		Status:     1,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	{
		TaskCore: TaskCore{
			Id:          6,
			Title:       "Submit expense report",
			Description: "Submit expenses for reimbursement",
			DueDate:     time.Now().AddDate(0, 0, 5).Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     3,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	// Tasks for the rest of the month
	{
		TaskCore: TaskCore{
			Id:          7,
			Title:       "Project kickoff meeting",
			Description: "Kick off the new project with the team",
			DueDate:     time.Now().AddDate(0, 0, 12).Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     1,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	{
		TaskCore: TaskCore{
			Id:          8,
			Title:       "Monthly report",
			Description: "Compile and send out the monthly progress report",
			Recurring:   "Monthly",
			DueDate:     time.Now().AddDate(0, 0, 13).Format("2006-01-02"),
		},
		Status:     1,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	{
		TaskCore: TaskCore{
			Id:          9,
			Title:       "Training session",
			Description: "Attend the training session on new software tools",
			DueDate:     time.Now().AddDate(0, 0, 14).Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     1,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	{
		TaskCore: TaskCore{
			Id:          10,
			Title:       "Website redesign",
			Description: "Discuss initial ideas for website redesign",
			DueDate:     time.Now().AddDate(0, 0, 16).Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     1,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	// 3 Overdue tasks
	{
		TaskCore: TaskCore{
			Id:          11,
			Title:       "Follow up on overdue invoice",
			Description: "Contact the client regarding the overdue invoice",
			DueDate:     time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     4,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	{
		TaskCore: TaskCore{
			Id:          12,
			Title:       "Bug fixing",
			Description: "Fix critical bugs reported by QA team",
			DueDate:     time.Now().AddDate(0, 0, -3).Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     4,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	{
		TaskCore: TaskCore{
			Id:          13,
			Title:       "Update documentation",
			Description: "Update project documentation with latest changes",
			DueDate:     time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     4,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	// 3 In Progress tasks
	{
		TaskCore: TaskCore{
			Id:          14,
			Title:       "Code refactoring",
			Description: "Refactor the legacy codebase for better performance",
			DueDate:     time.Now().Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     3,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	{
		TaskCore: TaskCore{
			Id:          15,
			Title:       "UI redesign",
			Description: "Work on redesigning the user interface",
			DueDate:     time.Now().Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     3,
		CreatedBy:  1,
		AssignedTo: 1,
	},
	{
		TaskCore: TaskCore{
			Id:          16,
			Title:       "Client demo",
			Description: "Prepare for the upcoming client demo",
			DueDate:     time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
			Recurring:   "false",
		},
		Status:     3,
		CreatedBy:  1,
		AssignedTo: 1,
	},
}

func Seed() error {
	fmt.Println("Seeding the db")
	db, err := sql.Open("sqlite3", "file:scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	db.Exec("DELETE from \"tasks\"")

	for i, task := range tasks {
		err := CreateTask(task)
		if err != nil {
			fmt.Printf("Task number %d failed: %s", i, err.Error())
		}
	}
	return nil
}
