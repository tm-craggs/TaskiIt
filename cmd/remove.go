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
	all      bool
	complete bool
	priority bool
	open     bool
	normal   bool
}

// Helper function to parse flags for removeCmd
func getRemoveFlags(cmd *cobra.Command) (*removeFlags, error) {
	flags := &removeFlags{}
	var err error

	if flags.all, err = cmd.Flags().GetBool("all"); err != nil {
		return nil, fmt.Errorf("failed to parse --all flag: %w", err)
	}
	if flags.complete, err = cmd.Flags().GetBool("complete"); err != nil {
		return nil, fmt.Errorf("failed to parse --complete flag: %w", err)
	}
	if flags.priority, err = cmd.Flags().GetBool("priority"); err != nil {
		return nil, fmt.Errorf("failed to parse --priority flag: %w", err)
	}
	if flags.open, err = cmd.Flags().GetBool("open"); err != nil {
		return nil, fmt.Errorf("failed to parse --open flag: %w", err)
	}
	if flags.normal, err = cmd.Flags().GetBool("normal"); err != nil {
		return nil, fmt.Errorf("failed to parse --normal flag: %w", err)
	}

	return flags, nil
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove all task",
	Long:  `Long description here`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		flags, err := getRemoveFlags(cmd)
		if err != nil {
			return err
		}

		// --all cannot be combined with other flags or args
		if flags.all {
			if flags.complete || flags.priority || flags.open || flags.normal || len(args) > 0 {
				return fmt.Errorf("--all cannot be combined with other flags or task IDs")
			}
		}

		noFlags := !flags.all && !flags.complete && !flags.priority && !flags.open && !flags.normal
		if noFlags && len(args) == 0 {
			return fmt.Errorf("task ID required")
		}

		// If all task ID is provided, check that it is valid
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
		case flags.all:
			confirmationText = "Remove all tasks?"
		case flags.complete:
			confirmationText = "Remove all completed tasks?"
		case flags.priority, flags.open, flags.normal:
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

		if flags.all {
			if err := task.RemoveAllTasks(); err != nil {
				return fmt.Errorf("failed to remove all tasks: %w", err)
			}
			return nil
		}

		if flags.complete || flags.priority || flags.open || flags.normal {
			tasks, err := task.GetTasks()
			if err != nil {
				return fmt.Errorf("failed to get tasks: %w", err)
			}

			filteredTasks := util.FilterTasks(tasks, flags.complete, flags.priority, flags.open, flags.normal)

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
	removeCmd.Flags().BoolP("all", "all", false, "Remove all tasks")
	removeCmd.Flags().BoolP("complete", "c", false, "Remove all complete tasks")
	removeCmd.Flags().BoolP("open", "o", false, "Remove all open tasks")
	removeCmd.Flags().BoolP("priority", "p", false, "Remove all priority tasks")
	removeCmd.Flags().BoolP("normal", "n", false, "Remove all normal priority tasks")

	rootCmd.AddCommand(removeCmd)
}
