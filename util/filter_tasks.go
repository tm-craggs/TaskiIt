package util

import "github.com/tcraggs/TidyTask/task"

func FilterTasks(tasks []task.Task,
	filterComplete bool, filterPriority bool, filterNotComplete bool, filterNotPriority bool) []task.Task {

	// put tasks through filters
	var filteredTasks []task.Task

	for _, t := range tasks {
		if filterPriority && !t.Priority {
			continue
		}
		if filterComplete && !t.Complete {
			continue
		}
		if filterNotComplete && t.Complete {
			continue
		}
		if filterNotPriority && t.Priority {
			continue
		}
		filteredTasks = append(filteredTasks, t)
	}

	return filteredTasks

}
