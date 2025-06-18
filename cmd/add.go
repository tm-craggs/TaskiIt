package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
	"strings"
)

// create struct that defines the available flags for add command
type addFlags struct {
	due      string
	priority bool
}

// helper function to parse flags with error handling
func getAddFlags(cmd *cobra.Command) (addFlags, error) {
	var flags addFlags
	var err error

	flags.due, err = cmd.Flags().GetString("due")
	if err != nil {
		return flags, fmt.Errorf("failed to parse --due flag: %w", err)
	}

	flags.priority, err = cmd.Flags().GetBool("priority")
	if err != nil {
		return flags, fmt.Errorf("failed to parse --priority flag: %w", err)
	}

	return flags, nil
}

// addCmd represents the add subcommand
var addCmd = &cobra.Command{
	Use:   "add [task name]... [flags]",
	Short: "Add a new task to your to-do list",
	Long:  `Add a new task to your to-do list. Optionally, use --due to add a due date, and --priority to mark task as high priority`,
	Example: `  tidytask add "Finish Homework"
    Add Finish Homework to your to-do list

  tidytask add Submit Essay --due 02-01-2006
	Add "Submit Essay" to your to-do list with 2nd of January 2006 as the due date

  tidytask add E-Mail boss --priority
    Add "E-Mail boss" to your to-do list and mark task as high priority

  tidytask add Finish Project --due 02-01-2006 --priority
    Add "Finish Project" to your to-do list with 2nd of January 2006 as the due date and mark task as high priority`,

	// backup database before running command
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := task.BackupDB()
		if err != nil {
			return fmt.Errorf("failed to back up database: %w", err)
		}
		return nil
	},

	// add command logic
	RunE: func(cmd *cobra.Command, args []string) error {

		// check if task name has been provided
		if len(args) == 0 {
			return fmt.Errorf("no task name provided")
		}

		// get flags
		flags, err := getAddFlags(cmd)
		if err != nil {
			return err
		}

		// if due date is provided, check it matches format DD-MM-YYYY
		if flags.due != "" {
			if err := util.VerifyDate(flags.due); err != nil {
				return fmt.Errorf("invalid date format. use DD-MM-YYYY")
			}
		}

		// join each args value to make task title
		title := strings.Join(args, " ")

		// create new task with input values
		newTask := task.Task{
			Title:    title,
			Due:      flags.due,
			Complete: false,
			Priority: flags.priority,
		}

		// add task to database
		if err := task.AddTask(newTask); err != nil {
			return fmt.Errorf("failed to add task: %w", err)
		}

		// exit
		fmt.Println("Task Added")
		return nil
	},
}

func init() {
	addCmd.Flags().StringP("due", "d", "", "Add a due date to task (DD-MM-YYYY)")
	addCmd.Flags().BoolP("priority", "p", false, "Mark task as high priority")

	rootCmd.AddCommand(addCmd)
}
