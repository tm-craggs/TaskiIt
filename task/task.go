package task

type Task struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Due      string `json:"due"`
	Complete bool   `json:"complete"`
	Priority bool   `json:"priority"`
}
