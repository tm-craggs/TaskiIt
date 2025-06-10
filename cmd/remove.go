package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
	"strconv"
)

// Struct to hold flag values for the remove command
type removeFlags struct {
	All      bool
	Complete bool
	Priority bool
	Open     bool
	Normal   bool
}

// Helper function to parse flags for removeCmd
func getRemoveFlags(cmd *cobra.Command) (*removeFlags, error) {
	f := &removeFlags{}
	var err error

	if f.All, err = cmd.Flags().GetBool("all"); err != nil {
		return nil, fmt.Errorf("failed to parse --all flag: %w", err)
	}
	if f.Complete, err = cmd.Flags().GetBool("complete"); err != nil {
		return nil, fmt.Errorf("failed to parse --complete flag: %w", err)
	}
	if f.Priority, err = cmd.Flags().GetBool("priority"); err != nil {
		return nil, fmt.Errorf("failed to parse --priority flag: %w", err)
	}
	if f.Open, err = cmd.Flags().GetBool("open"); err != nil {
		return nil, fmt.Errorf("failed to parse --open flag: %w", err)
	}
	if f.Normal, err = cmd.Flags().GetBool("normal"); err != nil {
		return nil, fmt.Errorf("failed to parse --normal flag: %w", err)
	}

	return f, nil
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a task",
	Long:  `Long description here`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		flags, err := getRemoveFlags(cmd)
		if err != nil {
			return err
		}

		// --all cannot be combined with other flags or args
		if flags.All {
			if flags.Complete || flags.Priority || flags.Open || flags.Normal || len(args) > 0 {
				return fmt.Errorf("--all cannot be combined with other flags or task IDs")
			}
		}

		noFlags := !flags.All && !flags.Complete && !flags.Priority && !flags.Open && !flags.Normal
		if noFlags && len(args) == 0 {
			return fmt.Errorf("task ID required")
		}

		// If a task ID is provided, check that it is valid
		if len(args) > 0 {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task ID: %s", args[0])
			}

			if err := task.CheckTaskExists(id); err != nil {
				return fmt.Errorf("task does not exist: %w", err)
			}
		}

		// Prepare confirmation text
		var confirmationText string
		switch {
		case flags.All:
			confirmationText = "Remove all tasks?"
		case flags.Complete:
			confirmationText = "Remove all completed tasks?"
		case flags.Priority, flags.Open, flags.Normal:
			confirmationText = "Remove filtered tasks?"
		default:
			confirmationText = "Remove task " + args[0] + "?"
		}

		if !util.ConfirmAction(confirmationText) {
			cmd.SilenceUsage = true
			return fmt.Errorf("aborted by user")
		}

		if err := task.BackupDB(); err != nil {
			return fmt.Errorf("failed to back up database: %w", err)
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := getRemoveFlags(cmd)
		if err != nil {
			return fmt.Errorf("error parsing flags: %w", err)
		}

		if flags.All {
			if err := task.RemoveAllTasks(); err != nil {
				return fmt.Errorf("failed to remove all tasks: %w", err)
			}
			return nil
		}

		if flags.Complete || flags.Priority || flags.Open || flags.Normal {
			tasks, err := task.GetTasks()
			if err != nil {
				return fmt.Errorf("failed to get tasks: %w", err)
			}

			filteredTasks := util.FilterTasks(tasks, flags.Complete, flags.Priority, flags.Open, flags.Normal)

			if err := task.RemoveInputTasks(filteredTasks); err != nil {
				return fmt.Errorf("failed to remove filtered tasks: %w", err)
			}
			return nil
		}

		// Remove single task by ID
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %s", args[0])
		}

		if err := task.RemoveTask(id); err != nil {
			return fmt.Errorf("failed to remove task %d: %w", id, err)
		}

		return nil
	},
}

func init() {
	removeCmd.Flags().BoolP("all", "a", false, "Remove all tasks")
	removeCmd.Flags().BoolP("complete", "c", false, "Remove all complete tasks")
	removeCmd.Flags().BoolP("open", "o", false, "Remove all open tasks")
	removeCmd.Flags().BoolP("priority", "p", false, "Remove all priority tasks")
	removeCmd.Flags().BoolP("normal", "n", false, "Remove all normal priority tasks")

	rootCmd.AddCommand(removeCmd)
}
