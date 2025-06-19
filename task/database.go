package task

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DB is a global variable representing the database connection
var DB *sql.DB

func getDBPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Failed to get filepath for config directory", err)
	}
	appDir := filepath.Join(configDir, "tidytask")

	// create app directory if it doesn't exist
	if err := os.MkdirAll(appDir, 0755); err != nil {
		log.Fatal("Failed to create config directory: ", err)
	}

	return filepath.Join(appDir, "tasks.db")
}

// InitDB creates the SQLite database by creating the database file and the tasks table if it does not exist
func InitDB() error {

	// open, or create if not exists, the SQLite database file "tasks.db"
	var err error

	// create database at dbPath
	dbPath := getDBPath()
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// define database table according to the Task struct
	createTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		due TEXT,
		complete BOOLEAN NOT NULL DEFAULT false,
		priority BOOLEAN NOT NULL DEFAULT false,
		complete_date TEXT          
	);`

	// execute the SQL statement, if error: log and exit
	_, err = DB.Exec(createTable)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// CloseDB safely closes the database connection, if it had been opened.
// it returns any error encountered during close, and if the DB was never initialised, it returns nil
// this function is always called on exit of tidytask
func CloseDB() error {
	// check if DB connection exists before trying to close it
	if DB != nil {
		// close the connection and return any errors
		return DB.Close()
	}
	// nothing to close, return nil
	return nil
}

// CheckTaskExists checks if a task with the given ID exists in the database.
// It returns an error if the task does not exist or if there is a query error.
func CheckTaskExists(id int) error {

	// exists will store whether the task exists (true) or not (false)
	var exists bool

	// SQL query returns true if at least one row with the given ID exists in the tasks table
	query := "SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)"

	// execute query and scan the result into the exists variable
	err := DB.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("query error checking task existence: %w", err)
	}

	// if exists is false, no task with ID was found, return an error
	if !exists {
		return fmt.Errorf("task with ID %d does not exist", id)
	}

	// if the task exists, return no error
	return nil
}

// AddTask inserts a new task into the tasks database table by taking in a task struct
func AddTask(t Task) error {

	// SQL insert statement to add a new task, setting complete_date to NULL
	stmt := `INSERT INTO tasks (title, due, complete, priority, complete_date) VALUES (?, ?, ?, ?, NULL)`

	// execute the insert statement with the task's fields as parameters
	_, err := DB.Exec(stmt, t.Title, t.Due, t.Complete, t.Priority)

	// return any error encountered
	return err
}

// RemoveTask deletes the task with the specified ID from the database.
func RemoveTask(id int) error {

	// execute DELETE SQL statement to remove the task matching the given ID
	_, err := DB.Exec("DELETE FROM tasks WHERE id = ?", id)

	// return any error encountered
	return err
}

// CompleteTask marks the task with the specified ID in the database as complete
func CompleteTask(id int) error {

	// get current date in layout YYYY-MM-DD
	currentDate := time.Now().Format("2006-01-02")

	// execute SQL statement to mark task as complete
	// only update complete_date if field is NULL
	_, err := DB.Exec(`
		UPDATE tasks
		SET complete = 1,
		    complete_date = CASE
		        WHEN complete_date IS NULL THEN ?
		        ELSE complete_date
		    END
		WHERE id = ?
	`, currentDate, id)

	// return any error encountered
	return err
}

// ReopenTask updates the task with the given ID to mark it as open (incomplete)
// it sets the 'complete' field to false and clears complete_date
func ReopenTask(id int) error {

	// execute UPDATE SQL statement to set complete to false and clear completion date
	_, err := DB.Exec("UPDATE tasks SET complete = 0, complete_date = NULL WHERE id = ?", id)

	// return any error encountered
	return err
}

// GetTasks retrieves all tasks from the database and returns them as a slice of Task structs.
// tasks are ordered by completion status, priority, presence of a due date, and due date ascending.
func GetTasks() ([]Task, error) {

	// SQL query to select all columns
	query := `
		SELECT id, title, due, complete, priority, complete_date
		FROM tasks
		ORDER BY
		    complete ASC, -- incomplete tasks first, ASC puts false (0) before true (1)
			priority DESC, -- among incomplete tasks, priority DESC puts priority tasks first
			due IS NOT NULL DESC,  -- tasks with a due date come before tasks without a due date
			due ASC  -- tasks are sorted by ascending due date, earliest first
	`

	// execute query and get each row
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	// close rows
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Failed to close rows: %v", err)
			return
		}
	}(rows)

	// slice to hold all retrieved tasks
	var tasks []Task

	// iterate over each row returned by query
	for rows.Next() {

		// create empty task struct
		var t Task

		// scan the columns of the current row into the Task struct
		err := rows.Scan(&t.ID, &t.Title, &t.Due, &t.Complete, &t.Priority, &t.CompleteDate)
		if err != nil {
			// return nil and error if scanning fails
			return nil, err
		}
		// add the populated task struct to the tasks slice
		tasks = append(tasks, t)
	}

	// return the slice of tasks, and return nil error on success
	return tasks, nil
}

// SetDue updates the due date of the task identified by the given ID.
// it sets the task's due field to the provided newDate string
func SetDue(id int, newDate string) error {

	// execute an UPDATE SQL statement to replace the due date for the specified task ID
	_, err := DB.Exec("UPDATE tasks SET due = ? WHERE id = ?", newDate, id)

	// return any encountered error
	return err
}

// SetTitle updates the due date of the task identified by the given ID.
// it sets the task's title field to the provided newTitle string
func SetTitle(id int, newTitle string) error {

	// execute an UPDATE SQL statement to replace the title for the specified task ID
	_, err := DB.Exec("UPDATE tasks SET title = ? WHERE id = ?", newTitle, id)

	// return any encountered error
	return err
}

// TogglePriority flips the priority status of the task identified by the given ID.
func TogglePriority(id int) error {

	// execute an UPDATE SQL statement to flip the priority status for the specified task ID
	_, err := DB.Exec("UPDATE tasks SET priority = NOT priority WHERE id = ?", id)

	// return any encountered error
	return err
}

// SearchTasks searches the tasks database for tasks where the given keyword matches any of the specified fields
// Returns a slice of matching Task structs or an error if the query fails
func SearchTasks(keyword string, searchID bool, searchTitle bool, searchDue bool) ([]Task, error) {

	// conditions holds individual SQL WHERE clauses for each enabled search field
	var conditions []string

	// args holds the arguments for the parameterised query placeholders
	var args []interface{}

	// if searching by ID, add a condition to match keyword against the ID cast as text
	if searchID {
		conditions = append(conditions, "CAST(id AS TEXT) LIKE ?")
		args = append(args, "%"+keyword+"%")
	}

	// if searching by title, add a LIKE condition for the title column
	if searchTitle {
		conditions = append(conditions, "title LIKE ?")
		args = append(args, "%"+keyword+"%")
	}

	// searching by due date, add a LIKE condition for the due column
	if searchDue {
		conditions = append(conditions, "due LIKE ?")
		args = append(args, "%"+keyword+"%")
	}

	// generate full query joining all conditions with OR
	query := `
		SELECT id, title, due, complete, priority
		FROM tasks
		WHERE ` + strings.Join(conditions, " OR ") + `
		ORDER BY complete ASC, priority DESC, due IS NOT NULL DESC, due ASC`

	// execute query and get rows, return nil with error if fails
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	// close rows
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Failed to close rows: %v", err)
			return
		}
	}(rows)

	// slice to hold matched tasks
	var tasks []Task

	// scan the columns of the current row into a Task struct
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Due, &t.Complete, &t.Priority)

		if err != nil {
			// return nil and error if scanning fails
			return nil, err
		}

		// add Task struct to tasks slice
		tasks = append(tasks, t)
	}

	// return task slice, and nil error on success
	return tasks, nil

}

// BackupDB creates a backup copy of the current SQLite database file ("tasks.db").
// It reads the original database file and writes its contents to a new file "tasks.db.bak".
func BackupDB() error {

	dbPath := getDBPath()
	backupPath := dbPath + ".bak"

	// read contents of current database file into memory
	input, err := os.ReadFile(dbPath)
	if err != nil {
		return err
	}

	// write the contents into the backup file with permission set to 0644
	return os.WriteFile(backupPath, input, 0644)
}

// RestoreBackup replaces the current database file with the backup copy.
// After successfully restoring, it deletes the backup file.
func RestoreBackup() error {

	// get path for DB
	dbPath := getDBPath()

	// get path for backup
	backupPath := dbPath + ".bak"

	// read contents of backup database file into memory
	input, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}

	// overwrite the contents of main database with backup contents
	err = os.WriteFile(dbPath, input, 0644)
	if err != nil {
		return err
	}

	// delete the backup file after successful restore
	err = os.Remove(backupPath)

	// log a warning if backup could not be deleted, but do not treat as failure
	if err != nil {
		// Not a fatal error, but you might want to log it
		fmt.Println("Warning: failed to delete backup file:", err)
	}

	return nil
}
