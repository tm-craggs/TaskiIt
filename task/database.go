package task

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
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
		due TEXT,
		complete BOOLEAN NOT NULL DEFAULT false,
		priority BOOLEAN NOT NULL DEFAULT false,
		complete_date TEXT          
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

func AddTask(t Task) error {
	stmt := `INSERT INTO tasks (title, due, complete, priority) VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(stmt, t.Title, t.Due, t.Complete, t.Priority)
	return err
}

func DeleteTask(id int) error {
	_, err := DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

func DeleteAllTasks() error {
	// remove all items from database
	_, err := DB.Exec("DELETE FROM tasks")
	if err != nil {
		return err
	}

	// reset auto-increment IDs to 1
	_, err = DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'tasks'")
	if err != nil {
		return err
	}

	return nil
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
		SELECT id, title, due, complete, priority
		FROM tasks
		ORDER BY
		    complete ASC,
			priority DESC,             -- High priority first (only applies to incomplete now)
			due IS NOT NULL DESC, -- Tasks with deadlines first
			due ASC
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Failed to close rows: %v", err)
			return
		}
	}(rows)

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Due, &t.Complete, &t.Priority)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func SetDue(id int, newDate string) error {
	_, err := DB.Exec("UPDATE tasks SET due = ? WHERE id = ?", newDate, id)
	return err
}

func SetTitle(id int, newTitle string) error {
	_, err := DB.Exec("UPDATE tasks SET title = ? WHERE id = ?", newTitle, id)
	return err
}

func TogglePriority(id int) error {
	_, err := DB.Exec("UPDATE tasks SET priority = NOT priority WHERE id = ?", id)
	return err
}

func SearchTasks(keyword string, searchID bool, searchTitle bool, searchDue bool) ([]Task, error) {
	var conditions []string
	var args []interface{}

	if searchID {
		conditions = append(conditions, "CAST(id AS TEXT) LIKE ?")
		args = append(args, "%"+keyword+"%")
	}
	if searchTitle {
		conditions = append(conditions, "title LIKE ?")
		args = append(args, "%"+keyword+"%")
	}
	if searchDue {
		conditions = append(conditions, "due LIKE ?")
		args = append(args, "%"+keyword+"%")
	}

	query := `
		SELECT id, title, due, complete, priority
		FROM tasks
		WHERE ` + strings.Join(conditions, " OR ") + `
		ORDER BY complete ASC, priority DESC, due IS NOT NULL DESC, due ASC`

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Failed to close rows: %v", err)
			return
		}
	}(rows)

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Due, &t.Complete, &t.Priority)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil

}

func BackupDB() error {
	input, err := os.ReadFile("tasks.db")
	if err != nil {
		return err
	}
	return os.WriteFile("tasks.db.bak", input, 0644)
}

func RestoreBackup() error {
	input, err := os.ReadFile("tasks.db.bak")
	if err != nil {
		return err
	}

	err = os.WriteFile("tasks.db", input, 0644)
	if err != nil {
		return err
	}

	// Delete the backup file after successful restore
	err = os.Remove("tasks.db.bak")
	if err != nil {
		// Not a fatal error, but you might want to log it
		fmt.Println("Warning: failed to delete backup file:", err)
	}

	return nil
}
