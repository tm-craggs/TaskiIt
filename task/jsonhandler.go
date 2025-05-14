package task

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadTasks() error {
	file, err := os.Open("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			Tasks = []Task{}
			return nil
		}
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	return json.NewDecoder(file).Decode(&Tasks)
}

func SaveTask(task Task) error {
	const filename = "tasks.json"

	var tasks []Task

	// Check if the file exists
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode existing tasks if file is not empty
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.Size() > 0 {
		if err := json.NewDecoder(file).Decode(&tasks); err != nil {
			return fmt.Errorf("error decoding existing tasks: %w", err)
		}
	}

	// Append the new task
	tasks = append(tasks, task)

	// Truncate the file and rewrite it
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tasks)
}
