package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"strconv"
)

type completeFlags struct {
	All bool
}

func getCompleteFlags(cmd *cobra.Command) (completeFlags, error) {
	var f completeFlags
	var err error

	f.All, err = cmd.Flags().GetBool("all")
	if err != nil {
		return f, fmt.Errorf("failed to parse --all flag: %w", err)
	}

	return f, nil
}

// completeCmd represented the `complete` command for marking tasks as done
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Complete a task",
	Long:  `Long goes here`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := task.BackupDB(); err != nil {
			return fmt.Errorf("failed to back up database: %w", err)
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := getCompleteFlags(cmd)
		if err != nil {
			return err
		}

		if flags.All {
			if len(args) > 0 {
				return fmt.Errorf("--all flag cannot be used with task IDs")
			}

			if err := task.CompleteAllTasks(); err != nil {
				return fmt.Errorf("failed to complete all tasks: %w", err)
			}

			fmt.Println("All tasks marked complete")
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("task ID required")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("failed to parse task ID: %w", err)
		}

		if err := task.CheckTaskExists(id); err != nil {
			return fmt.Errorf("task does not exist: %w", err)
		}

		if err := task.CompleteTask(id); err != nil {
			return fmt.Errorf("failed to complete task: %w", err)
		}

		fmt.Println("Task completed successfully")
		return nil
	},
}

func init() {
	completeCmd.Flags().BoolP("all", "a", false, "Complete all tasks")
	rootCmd.AddCommand(completeCmd)
}
