package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tcraggs/TidyTask/task"
)

var (
	due      string
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
			return
		}

		// task title is taken in as args
		title := args[0]

		newTask := task.Task{
			Title:    title,
			Due:      due,
			Complete: false,
			Priority: *priority,
		}

		if err := task.AddTask(newTask); err != nil {
			fmt.Println("Failed to add task", err.Error())
			return
		}

		fmt.Println("Task Created")
		fmt.Println("ID:", newTask.ID)
		fmt.Println("Title:", newTask.Title)
		fmt.Println("Due:", newTask.Due)
		fmt.Println("Priority:", newTask.Priority)

	},
}

func init() {
	addCmd.Flags().StringVarP(&due, "due", "d", "", "Set a due (e.g. 2025-05-14)")
	addCmd.Flags().BoolVarP(priority, "priority", "p", false, "Mark the task as high priority")

	rootCmd.AddCommand(addCmd)

}
