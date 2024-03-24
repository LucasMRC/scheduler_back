package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jomei/notionapi"
	_ "github.com/mattn/go-sqlite3"
)

var notion NotionAPI

const createTasksTable string = `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER NOT NULL PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_DATE,
		description TEXT,
		title TEXT NOT NULL,
		due_date DATETIME NOT NULL,
		recurring INTEGER NOT NULL DEFAULT 1,
		assigned_to INTEGER NOT NULL ON CONFLICT FAIL,
		status INTEGER NOT NULL DEFAULT 1,
		created_by INTEGER,
		FOREIGN KEY (recurring) REFERENCES recurring_types(id),
		FOREIGN KEY (assigned_to) REFERENCES users(id),
		FOREIGN KEY (status) REFERENCES task_status(id),
		FOREIGN KEY (created_by) REFERENCES users(id)
	);`

const createUsersTable string = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_DATE,
		alias TEXT NOT NULL,
		email TEXT NOT NULL,
		hash TEXT NOT NULL
	);`

const createRecurringTypesTable string = `
	CREATE TABLE IF NOT EXISTS recurring_types (
		id INTEGER NOT NULL PRIMARY KEY,
		title TEXT NOT NULL UNIQUE
	);`
const insertRecurringTypes string = `
	INSERT OR IGNORE INTO recurring_types (title) VALUES
	('Once'),
	('Daily'),
	('Weekly'),
	('Monthly'),
	('Yearly');`

const createTaskStatusTable string = `
	CREATE TABLE IF NOT EXISTS task_status (
		id INTEGER NOT NULL PRIMARY KEY,
		title TEXT NOT NULL UNIQUE
	);`
const insertTaskStatus string = `
	INSERT OR IGNORE INTO task_status (title) VALUES
	('Pending'),
	('In Progress'),
	('Completed'),
	('Overdue');`

func CreateTables(db *sql.DB) error {
	fmt.Println("Creating recurring types table")
	_, err := db.Exec(createRecurringTypesTable)
	if err != nil {
		return err
	}
	fmt.Println("Creating task status table")
	_, err = db.Exec(createTaskStatusTable)
	if err != nil {
		return err
	}
	fmt.Println("Creating users table")
	_, err = db.Exec(createUsersTable)
	if err != nil {
		return err
	}
	fmt.Println("Creating tasks table")
	_, err = db.Exec(createTasksTable)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	fmt.Println("Initializing database")
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		panic(err)
	}

	fmt.Println("Creating tables")
	if err := CreateTables(db); err != nil {
		panic(err)
	}

	NOTION_DB_TOKEN := notionapi.Token(os.Getenv("DB_TOKEN"))
	notion = NotionAPI{
		Client: notionapi.NewClient(NOTION_DB_TOKEN),
	}
}

func SaveUser(user User) error {
	db, err := sql.Open("sqlite3", "./pkg/database/scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users (alias, email, hash) VALUES (?, ?, ?)", user.Alias, user.Email, user.Password)
	if err != nil {
		return err
	}
	go func() {
		notion.SaveUser(user)
	}()

	return nil
}

func GetUsers() ([]User, error) {
	db, err := sql.Open("sqlite3", "./pkg/database/scheduler.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT alias, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Alias, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetUser(alias string) (User, error) {
	db, err := sql.Open("sqlite3", "./pkg/database/scheduler.db")
	if err != nil {
		return User{}, err
	}
	defer db.Close()

	var user User
	err = db.QueryRow("SELECT alias, email FROM users WHERE alias = ?", alias).Scan(&user.Alias, &user.Email)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func CreateTask(task Task) error {
	db, err := sql.Open("sqlite3", "./pkg/database/scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO tasks (title, description, due_date, recurring, assigned_to) VALUES (?, ?, ?, ?, ?)", task.Title, task.Description, task.DueDate, task)
	if err != nil {
		return err
	}
	return nil
}

func GetTasks(alias string) ([]Task, error) {
	db, err := sql.Open("sqlite3", "./pkg/database/scheduler.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT title, description, due_date, recurring, assigned_to FROM tasks JOIN users ON assign_to = users.id WHERE users.alias = ?", alias)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Title, &task.Description, &task.DueDate, &task.Recurring, &task.AssignedTo)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetTask(id int) (Task, error) {
	db, err := sql.Open("sqlite3", "./pkg/database/scheduler.db")
	if err != nil {
		return Task{}, err
	}
	defer db.Close()

	var task Task
	err = db.QueryRow("SELECT title, description, due_date, recurring, assigned_to FROM tasks WHERE id = ?", id).Scan(&task.Title, &task.Description, &task.DueDate, &task)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func UpdateTask(taskId string, task Task) error {
	db, err := sql.Open("sqlite3", "./pkg/database/scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE tasks SET title = ?, description = ?, due_date = ?, recurring = ?, assigned_to = ? WHERE id = ?", task.Title, task.Description, task.DueDate, taskId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTask(id string) error {
	db, err := sql.Open("sqlite3", "./pkg/database/scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
