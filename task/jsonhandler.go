package task

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
