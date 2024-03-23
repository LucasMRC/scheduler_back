package database

import (
	"database/sql"
	"fmt"

	// "github.com/LucasMRC/lb_back/pkg/notion"
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
		// notion.SaveUser(user)
	}()

	return nil
}
