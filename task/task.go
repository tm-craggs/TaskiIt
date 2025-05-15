package task

type Task struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Deadline string `json:"deadline"`
	Complete bool   `json:"complete"`
	Priority bool   `json:"priority"`
}
