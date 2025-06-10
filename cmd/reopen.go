package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"strconv"
)

type reopenFlags struct {
	all bool
}

func getReopenFlags(cmd *cobra.Command) (reopenFlags, error) {
	var flags reopenFlags
	var err error

	flags.all, err = cmd.Flags().GetBool("all")
	if err != nil {
		return flags, fmt.Errorf("failed to parse --all flag: %w", err)
	}

	return flags, nil
}

var reopenCmd = &cobra.Command{
	Use:   "reopen",
	Short: "Reopen all completed task",
	Long:  `Long goes here`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := task.BackupDB(); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := getReopenFlags(cmd)
		if err != nil {
			return err
		}

		if flags.all {
			if len(args) > 0 {
				return fmt.Errorf("use --all only to reopen all tasks")
			}

			if err := task.ReopenAllTasks(); err != nil {
				return fmt.Errorf("failed to reopen all tasks: %w", err)
			}

			fmt.Println("all tasks reopened")
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("please specify the task ID to reopen")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %w", err)
		}

		if err := task.CheckTaskExists(id); err != nil {
			return fmt.Errorf("task does not exist: %w", err)
		}

		if err := task.ReopenTask(id); err != nil {
			return fmt.Errorf("failed to reopen task: %w", err)
		}

		fmt.Println("Task reopened")
		return nil
	},
}

func init() {
	reopenCmd.Flags().BoolP("all", "all", false, "Reopen all tasks")
	rootCmd.AddCommand(reopenCmd)
}
