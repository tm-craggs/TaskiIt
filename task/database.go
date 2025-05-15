package task

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		deadline TEXT,
		complete BOOLEAN NOT NULL DEFAULT false,
		priority BOOLEAN NOT NULL DEFAULT false
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

func GetAllTasks() ([]Task, error) {
	rows, err := DB.Query("SELECT id, title, deadline, complete, priority FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Deadline, &t.Complete, &t.Priority)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func AddTask(t Task) error {
	stmt := `INSERT INTO tasks (title, deadline, complete, priority) VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(stmt, t.Title, t.Deadline, t.Complete, t.Priority)
	return err
}

func DeleteTask(id int) error {
	_, err := DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

func DeleteAllTasks() error {
	_, err := DB.Exec("DELETE FROM tasks")
	return err
}

func CompleteTask(id int) error {
	_, err := DB.Exec("UPDATE tasks SET complete = 1 WHERE id = ?", id)
	return err
}

func CompleteAllTasks() error {
	_, err := DB.Exec("UPDATE tasks SET complete = 1")
	return err
}

func ReopenTask(id int) error {
	_, err := DB.Exec("UPDATE tasks SET complete = 0 WHERE id = ?", id)
	return err
}

func ReopenAllTasks() error {
	_, err := DB.Exec("UPDATE tasks SET complete = 0")
	return err
}

func GetTasks() ([]Task, error) {
	query := `
		SELECT id, title, deadline, complete, priority
        FROM tasks
		ORDER BY
			priority DESC,
			deadline IS NOT NOT NULL DESC,
			deadline ASC `

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Deadline, &t.Complete, &t.Priority)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil

}
