package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"strconv"
)

// completeCmd represented the `complete` command for marking tasks as done
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Complete a task",
	Long:  `Long goes here`,

	// PreRunE runs before the main logic
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := task.BackupDB()

		if err != nil {
			// return error to PreRunE
			return fmt.Errorf("failed to back up database: %w", err)
		}

		// no error, continue to main logic
		return nil
	},

	// core logic for `complete` command
	RunE: func(cmd *cobra.Command, args []string) error {

		// check if --all flag set
		var completeAll, err = cmd.Flags().GetBool("all")
		if err != nil {
			return fmt.Errorf("failed to get flag options: %w", err)
		}

		// code for --all flag
		if completeAll {

			if len(args) > 0 {
				return fmt.Errorf("--all flag cannot be used with task IDs")
			}

			if err := task.CompleteAllTasks(); err != nil {
				return fmt.Errorf("failed to complete all tasks: %w", err)
			}

			// print confirmation and return
			fmt.Println("All tasks marked complete")
			return nil
		}

		// single task removal for single logic

		// check input has been given
		if len(args) == 0 {
			return fmt.Errorf("task ID required")
		}

		// convert args input to int
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("failed to parse task ID: %w", err)
		}

		// check if input ID exists in DB
		exists, err := task.CheckTaskExists(id)
		if err != nil {
			return fmt.Errorf("failed to check task existence: %w", err)
		}

		// if not exists, throw error
		if !exists {
			return fmt.Errorf("task ID not found")
		}

		// call complete task
		if err := task.CompleteTask(id); err != nil {
			return fmt.Errorf("failed to complete task: %w", err)
		}

		// print confirmation
		fmt.Println("Task completed successfully")
		return nil
	},
}

func init() {
	completeCmd.Flags().BoolP("all", "a", false, "Complete all tasks")
	rootCmd.AddCommand(completeCmd)
}
