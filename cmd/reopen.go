package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"strconv"
)

var reopenAll bool

var reopenCmd = &cobra.Command{
	Use:   "reopen",
	Short: "Reopen a completed task",
	Long:  `Long goes here`,
	Run: func(cmd *cobra.Command, args []string) {

		// code for --all flag
		if reopenAll {

			if len(args) > 0 {
				fmt.Println("Use --all only to reopen all tasks")
				return
			}

			if err := task.ReopenAllTasks(); err != nil {
				fmt.Println("Failed to reopen all tasks", err.Error())
				return
			}

			fmt.Println("All tasks reponed")

		} else {

			// code for single removal

			if len(args) == 0 {
				fmt.Println("Please specify the task ID to reopen")
				return
			}

			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Invalid task ID")
			}

			if err := task.ReopenTask(id); err != nil {
				fmt.Println("Failed to reopen task:", err.Error())
				return
			}

			fmt.Println("Task reponed")

		}
	},
}

func init() {
	reopenCmd.Flags().BoolVarP(&completeAll, "all", "a", false, "Reopen all tasks")
	rootCmd.AddCommand(reopenCmd)
}
