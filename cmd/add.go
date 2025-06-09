package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
	"strings"
)

// addCmd represents the add subcommand
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new task to your to-do list",
	Long:  `Long description goes here`,

	// PreRunE runs before main logic
	PreRunE: func(cmd *cobra.Command, args []string) error {

		err := task.BackupDB()

		if err != nil {
			// return error to PreRunE and print to standard output
			return fmt.Errorf("failed to back up database: %w", err)
		}

		// no error, continue to main logic
		return nil
	},

	// core logic for `add` command
	RunE: func(cmd *cobra.Command, args []string) error {

		// check if a task name was provided
		if len(args) == 0 {
			return fmt.Errorf("no task name provided")
		}

		// join all args into a single string to allow task names with spaces
		title := strings.Join(args, " ")

		// parse optional flags --due and --priority
		due, _ := cmd.Flags().GetString("due")
		priority, _ := cmd.Flags().GetBool("priority")

		// if a due date was provided, ensure it uses valid formatting
		if due != "" {
			if err := util.VerifyDate(due); err != nil {
				return fmt.Errorf("invalid date format. use YYYY-MM-DD")
			}
		}

		// create new task struct with provided values
		newTask := task.Task{
			Title:    title,
			Due:      due,
			Complete: false,
			Priority: priority,
		}

		// add the task to database using addTask function
		if err := task.AddTask(newTask); err != nil {
			return fmt.Errorf("failed to add task: %w", err)
		}

		// print confirmation
		fmt.Println("Task Created")
		return nil
	},
}

// init is called automatically to set up flags and register command
func init() {

	// define --due (-d) flag for setting a due date
	addCmd.Flags().StringP("due", "d", "", "Set a due (e.g. 2025-05-14)")

	// define the --priority (-p) flag for marking the task as high priority
	addCmd.Flags().BoolP("priority", "p", false, "Mark the task as high priority")

	// register the command with root so it is available to the CLI
	rootCmd.AddCommand(addCmd)

}
