package util

import "github.com/tm-craggs/tidytask/task"

// FilterTasks takes a slice of Task structs and returns a new slice
// only tasks that satisfy all enabled filters are added to the new slice
func FilterTasks(tasks []task.Task,
	filterComplete bool, filterPriority bool, filterNotComplete bool, filterNotPriority bool) []task.Task {

	// create slice to store tasks that pass all filters
	var filteredTasks []task.Task

	// iterate over each task in the input slice
	for _, t := range tasks {

		// skip task if filterPriority is active and task is not priority
		if filterPriority && !t.Priority {
			continue
		}

		// skip task if filterComplete is active and task is not complete
		if filterComplete && !t.Complete {
			continue
		}

		// skip the task if filterNotComplete is active and task is complete
		if filterNotComplete && t.Complete {
			continue
		}

		// skip task if filterNotPriority is active and task is priority
		if filterNotPriority && t.Priority {
			continue
		}

		// if task passed all filters, add to slice
		filteredTasks = append(filteredTasks, t)
	}

	// return filtered list of tasks
	return filteredTasks

}
