package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
	"strconv"
)

type editFlags struct {
	Title           string
	Due             string
	Priority        bool
	TitleChanged    bool
	DueChanged      bool
	PriorityChanged bool
}

func getEditFlags(cmd *cobra.Command) (editFlags, error) {
	var f editFlags
	var err error

	f.Title, err = cmd.Flags().GetString("title")
	if err != nil {
		return f, fmt.Errorf("failed to parse --title flag: %w", err)
	}

	f.Due, err = cmd.Flags().GetString("due")
	if err != nil {
		return f, fmt.Errorf("failed to parse --due flag: %w", err)
	}

	f.Priority, err = cmd.Flags().GetBool("priority")
	if err != nil {
		return f, fmt.Errorf("failed to parse --priority flag: %w", err)
	}

	f.TitleChanged = cmd.Flags().Changed("title")
	f.DueChanged = cmd.Flags().Changed("due")
	f.PriorityChanged = cmd.Flags().Changed("priority")

	return f, nil
}

// addCmd represents the add command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "add a new task to your to-do list",
	Long:  `Long description goes here`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := task.BackupDB(); err != nil {
			return err
		}

		if !util.ConfirmAction("Confirm Edit?") {
			cmd.SilenceUsage = true
			return fmt.Errorf("aborted by user")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return fmt.Errorf("no task ID provided")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID %w", err)
		}

		if err := task.CheckTaskExists(id); err != nil {
			return fmt.Errorf("task does not exist %w", err)
		}

		flags, err := getEditFlags(cmd)
		if err != nil {
			return err
		}

		if flags.TitleChanged {
			if err := task.SetTitle(id, flags.Title); err != nil {
				return fmt.Errorf("failed to update title: %w", err)
			}
		}

		if flags.PriorityChanged {
			if err := task.TogglePriority(id); err != nil {
				return fmt.Errorf("failed to toggle priority: %w", err)
			}
		}

		if flags.DueChanged {
			if err := util.VerifyDate(flags.Due); err != nil {
				return fmt.Errorf("invalid due date: %w", err)
			}
			if err := task.SetDue(id, flags.Due); err != nil {
				return fmt.Errorf("failed to update due date: %w", err)
			}
		}

		fmt.Println("Task updated")

		return nil
	},
}

func init() {
	editCmd.Flags().StringP("due", "d", "", "Set a due date (e.g. 2025-05-14)")
	editCmd.Flags().BoolP("priority", "p", false, "Mark the task as high priority")
	editCmd.Flags().StringP("title", "t", "", "Set a task title")

	rootCmd.AddCommand(editCmd)
}
