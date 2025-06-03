package task

import "database/sql"

type Task struct {
	ID           int            `json:"id"`
	Title        string         `json:"title"`
	Due          string         `json:"due"`
	Complete     bool           `json:"complete"`
	CompleteDate sql.NullString `json:"complete_date"`
	Priority     bool           `json:"priority"`
}
