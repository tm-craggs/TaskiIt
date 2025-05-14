package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tcraggs/TidyTask/task"
)

var (
	deadline string
	priority = pflag.BoolP("priority", "p", false, "Mark task as high priority")
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new task to your to-do list",
	Long:  `Long description goes here`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please enter a task name")
		}

		// task title is taken in as args
		title := args[0]

		// load current tasks in json into Task
		if err := task.LoadTasks(); err != nil {
			fmt.Printf("Error loading tasks: %s\n", err.Error())
			return
		}

		// find max ID
		nextID := 1

		// largest ID will typically be at the end of the list, so reverse traverse list
		for i := len(task.Tasks) - 1; i >= 0; i-- {
			t := task.Tasks[i]
			if t.ID >= nextID {
				nextID = t.ID + 1
				break // done, found max
			}
		}

		newTask := task.Task{
			ID:       nextID,
			Title:    title,
			Deadline: deadline,
			Complete: false,
			Priority: *priority,
		}

		task.SaveTask(newTask)

	},
}

func init() {
	addCmd.Flags().StringVarP(&deadline, "deadline", "d", "", "Set a deadline (e.g. 2025-05-14)")
	addCmd.Flags().BoolVarP(priority, "priority", "p", false, "Mark the task as high priority")

	rootCmd.AddCommand(addCmd)

}
