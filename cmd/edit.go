package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
	"strconv"
)

// create struct that defines the available flags for edit command
type editFlags struct {
	title           string
	due             string
	priority        bool
	titleChanged    bool
	dueChanged      bool
	priorityChanged bool
}

// helper function to parse flags with error handling
func getEditFlags(cmd *cobra.Command) (editFlags, error) {
	var flags editFlags
	var err error

	flags.title, err = cmd.Flags().GetString("title")
	if err != nil {
		return flags, fmt.Errorf("failed to parse --title flag: %w", err)
	}

	flags.due, err = cmd.Flags().GetString("due")
	if err != nil {
		return flags, fmt.Errorf("failed to parse --due flag: %w", err)
	}

	flags.priority, err = cmd.Flags().GetBool("priority")
	if err != nil {
		return flags, fmt.Errorf("failed to parse --priority flag: %w", err)
	}

	flags.titleChanged = cmd.Flags().Changed("title")
	flags.dueChanged = cmd.Flags().Changed("due")
	flags.priorityChanged = cmd.Flags().Changed("priority")

	return flags, nil
}

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [ID] [flags]",
	Short: "Edit details of an existing task",
	Long: `The 'edit' command allows you to update details of an existing task.

Specify the task ID and pass flags for the the details you wish to change.`,

	Example: `  tidytask edit 1 --title "Buy Groceries"
	Change the title of task 1 to "Buy Groceries"

  tidytask edit 3 --due 2006-01-02 --priority
	Change the due date of task 3 to 2nd of January 2006, and toggle the priority status

  tidytask edit 5 --title "Clean Room" --due 02-01-2006
	Change the title of task 5 to Clean Room and change the due date to 2nd of January 2006`,

	// main command logic
	RunE: func(cmd *cobra.Command, args []string) error {

		// check args
		if len(args) == 0 {
			return fmt.Errorf("no arguments provided; task ID required")
		}

		if len(args) > 1 {
			return fmt.Errorf("accepts 1 argument, received %d; use quotes for multi-word input", len(args))
		}

		// convert task ID to int
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID %w", err)
		}

		// check task exists
		if err := task.CheckTaskExists(id); err != nil {
			return fmt.Errorf("task does not exist %w", err)
		}

		// get flags
		flags, err := getEditFlags(cmd)
		if err != nil {
			return err
		}

		// backup database
		if err := task.BackupDB(); err != nil {
			fmt.Printf("Warning: failed to back up database: %v", err)
		}

		if !util.ConfirmAction("Confirm edit?") {
			return fmt.Errorf("aborted by user")
		}

		// update task title if title flagged
		if flags.titleChanged {
			if err := task.SetTitle(id, flags.title); err != nil {
				return fmt.Errorf("failed to update title: %w", err)
			}
		}

		// toggle task priority if priority flagged
		if flags.priorityChanged {
			if err := task.TogglePriority(id); err != nil {
				return fmt.Errorf("failed to toggle priority: %w", err)
			}
		}

		// update due date if due flagged
		if flags.dueChanged {
			if err := util.VerifyDate(flags.due); err != nil {
				return fmt.Errorf("invalid due date: %w", err)
			}
			if err := task.SetDue(id, flags.due); err != nil {
				return fmt.Errorf("failed to update due date: %w", err)
			}
		}

		// exit
		fmt.Println("Task updated")
		return nil
	},
}

// command initialisation
func init() {

	// define flags and add subcommand to root

	editCmd.Flags().StringP("due", "d", "", "Change due date of task (YYYY-MM-DD)")
	editCmd.Flags().BoolP("priority", "p", false, "Toggle the task priority")
	editCmd.Flags().StringP("title", "t", "", "Change the title of task")

	rootCmd.AddCommand(editCmd)
}
