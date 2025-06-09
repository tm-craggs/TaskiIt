package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
	"strconv"
)

// addCmd represents the add command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "add a new task to your to-do list",
	Long:  `Long description goes here`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := task.BackupDB()
		if err != nil {
			return err
		}

		if !util.ConfirmAction("Confirm Edit?") {
			cmd.SilenceUsage = true
			return fmt.Errorf("aborted by user")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		// check if a task name was provided
		if len(args) == 0 {
			return fmt.Errorf("no task ID provided")
		}

		// parse task ID
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID %w", err)
		}

		// check if task exists
		if err := task.CheckTaskExists(id); err != nil {
			return fmt.Errorf("task does not exist %w", err)
		}

		// get string following --title flag
		editTitle, err := cmd.Flags().GetString("title")
		if err != nil {
			return err
		}

		// get string following --due flag
		editDue, err := cmd.Flags().GetString("due")
		if err != nil {
			return err
		}

		// check if flag was used, store as boolean
		priorityChanged := cmd.Flags().Changed("priority")
		titleChanged := cmd.Flags().Changed("title")
		dueChanged := cmd.Flags().Changed("due")

		// apply edits if booleans are true

		if titleChanged {
			if err := task.SetTitle(id, editTitle); err != nil {
				return fmt.Errorf("failed to update title: %w", err)
			}
		}

		if priorityChanged {
			if err := task.TogglePriority(id); err != nil {
				return fmt.Errorf("failed to toggle priority: %w", err)
			}
		}

		if dueChanged {
			if err := util.VerifyDate(editDue); err != nil {
				return fmt.Errorf("invalid due date: %w", err)
			}
			if err := task.SetDue(id, editDue); err != nil {
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
