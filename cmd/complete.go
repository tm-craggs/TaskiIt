package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"strconv"
)

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Complete a task",
	Long:  `Long goes here`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := task.BackupDB()
		if err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		var completeAll, err = cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println(err)
			return
		}

		// code for --all flag
		if completeAll {

			if len(args) > 0 {
				fmt.Println("Use --all only to complete all tasks")
				return
			}

			if err := task.CompleteAllTasks(); err != nil {
				fmt.Println("Failed to complete all tasks", err.Error())
				return
			}

			fmt.Println("All tasks successfully")

		} else {

			// code for single removal

			if len(args) == 0 {
				fmt.Println("Please specify the task ID to remove")
				return
			}

			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Invalid task ID")
			}

			if err := task.CompleteTask(id); err != nil {
				fmt.Println("Failed to complete task:", err.Error())
				return
			}

			fmt.Println("Task completed successfully")

		}
	},
}

func init() {
	completeCmd.Flags().BoolP("all", "a", false, "Complete all tasks")
	rootCmd.AddCommand(completeCmd)
}
