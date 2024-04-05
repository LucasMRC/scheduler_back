package database

import (
	"database/sql"
	"fmt"

	// "os"

	_ "github.com/mattn/go-sqlite3"
)

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
		alias TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
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
	fmt.Println("Inserting recurring types")
	_, err = db.Exec(insertRecurringTypes)
	if err != nil {
		return err
	}
	fmt.Println("Creating task status table")
	_, err = db.Exec(createTaskStatusTable)
	if err != nil {
		return err
	}
	fmt.Println("Inserting task status")
	_, err = db.Exec(insertTaskStatus)
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
}

func SaveUser(user User) error {
	db, err := sql.Open("sqlite3", "file:scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users (alias, email, hash) VALUES (?, ?, ?)", user.Alias, user.Email, user.Hash)
	if err != nil {
		return err
	}
	// go func() {
	// 	NOTION_DB_TOKEN := notionapi.Token(os.Getenv("DB_TOKEN"))
	// 	notion := NotionAPI{
	// 		Client: notionapi.NewClient(NOTION_DB_TOKEN),
	// 	}
	// 	notion.SaveUser(user)
	// }()

	return nil
}

func GetUsers() ([]User, error) {
	db, err := sql.Open("sqlite3", "file:scheduler.db")
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
	db, err := sql.Open("sqlite3", "file:scheduler.db")
	if err != nil {
		return User{}, err
	}
	defer db.Close()

	var user User
	err = db.QueryRow("SELECT id, alias, email, hash FROM users WHERE alias = ?", alias).Scan(&user.Id, &user.Alias, &user.Email, &user.Hash)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func UpdateUser(alias string, user User) error {
	db, err := sql.Open("sqlite3", "file:scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE users SET email = ?, hash = ? WHERE alias = ?", user.Email, user.Hash, alias)
	if err != nil {
		return err
	}
	return nil
}

func CreateTask(task Task) error {
	db, err := sql.Open("sqlite3", "file:scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO tasks (title, description, due_date, recurring, assigned_to, created_by) VALUES (?, ?, ?, ?, ?, ?)", task.Title, task.Description, task.DueDate, task.Recurring, task.AssignedTo, task.CreatedBy)
	if err != nil {
		return err
	}
	// go func() {
	// 	NOTION_DB_TOKEN := notionapi.Token(os.Getenv("DB_TOKEN"))
	// 	notion := NotionAPI{
	// 		Client: notionapi.NewClient(NOTION_DB_TOKEN),
	// 	}
	// 	notion.CreateTask(task)
	// }()
	return nil
}

func GetTasks(alias string) ([]TaskDTO, error) {
	db, err := sql.Open("sqlite3", "file:scheduler.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT tasks.id, tasks.title, description, due_date, recurring, users.alias as assigned_to, task_status.title as status FROM tasks JOIN users ON tasks.assigned_to = users.id JOIN task_status on tasks.status = task_status.id WHERE users.alias = ?", alias)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]TaskDTO, 0)
	for rows.Next() {
		var task TaskDTO
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.DueDate, &task.Recurring, &task.AssignedTo, &task.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetTask(id int) (Task, error) {
	db, err := sql.Open("sqlite3", "file:scheduler.db")
	if err != nil {
		return Task{}, err
	}
	defer db.Close()

	var task Task
	err = db.QueryRow("SELECT title, description, due_date, recurring, assigned_to FROM tasks WHERE id = ?", id).Scan(&task.Title, &task.Description, &task.DueDate, &task.Recurring, &task.AssignedTo)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func UpdateTask(taskId int, task Task) (Task, error) {
	db, err := sql.Open("sqlite3", "file:scheduler.db")
	if err != nil {
		return Task{}, err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE tasks SET title = ?, description = ?, due_date = ?, recurring = ?, assigned_to = ? WHERE id = ?", task.Title, task.Description, task.DueDate, task.Recurring, task.AssignedTo, taskId)
	if err != nil {
		return Task{}, err
	}

	updatedTask, err := GetTask(taskId)
	if err != nil {
		return Task{}, err
	}
	// go func() {
	// 	NOTION_DB_TOKEN := notionapi.Token(os.Getenv("DB_TOKEN"))
	// 	notion := NotionAPI{
	// 		Client: notionapi.NewClient(NOTION_DB_TOKEN),
	// 	}
	// 	notion.UpdateTask(taskId, updatedTask)
	// }()
	return updatedTask, nil
}

func DeleteTask(id int) error {
	db, err := sql.Open("sqlite3", "file:scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	// go func() {
	// 	NOTION_DB_TOKEN := notionapi.Token(os.Getenv("DB_TOKEN"))
	// 	notion := NotionAPI{
	// 		Client: notionapi.NewClient(NOTION_DB_TOKEN),
	// 	}
	// 	notion.DeleteTask(id)
	// }()
	return nil
}
