package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"strconv"
)

var removeAll bool

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a task",
	Long:  `Long goes here`,
	Run: func(cmd *cobra.Command, args []string) {

		// code for --all flag
		if removeAll {

			if len(args) > 0 {
				fmt.Println("Use --all only to remove all tasks")
				return
			}

			if err := task.DeleteAllTasks(); err != nil {
				fmt.Println("Failed to delete all tasks", err.Error())
				return
			}

			fmt.Println("All tasks removed successfully")

		} else {

			// code for single removal

			if len(args) == 0 {
				fmt.Println("Please specify the task ID to remove")
				return
			}

			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Task ID not found")
			}

			if err := task.DeleteTask(id); err != nil {
				fmt.Println("Failed to remove task:", err.Error())
				return
			}

			fmt.Println("Task removed successfully")

		}
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&removeAll, "all", "a", false, "Remove all tasks")
	rootCmd.AddCommand(removeCmd)
}
