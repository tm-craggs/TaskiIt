package task

import "database/sql"

// Task represents a to-do list task
type Task struct {
	ID           int            `json:"id"`            // Unique ID for task (primary key)
	Title        string         `json:"title"`         // Title or description of the task (mandatory)
	Due          string         `json:"due"`           // Due date as string (empty string represents no due set)
	Complete     bool           `json:"complete"`      // Flag indicating the tasks completion status
	CompleteDate sql.NullString `json:"complete_date"` // Nullable date string representing when task was completed
	Priority     bool           `json:"priority"`      // Flag indicating if the task is marked as high priority
}
