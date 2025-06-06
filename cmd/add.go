package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
	"os"
	"strings"
)

// addCmd represents the add subcommand
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new task to your to-do list",
	Long:  `Long description goes here`,

	// PreRunE runs before main logic
	PreRunE: func(cmd *cobra.Command, args []string) error {

		// backup the database before making changes
		err := task.BackupDB()
		if err != nil {
			return err
		} // return error if backup fails

		return nil // otherwise, continue
	},

	// core logic for `add` command
	Run: func(cmd *cobra.Command, args []string) {

		// check if a task name was provided
		if len(args) == 0 {
			fmt.Println("Error: Task name is required")
			os.Exit(1)
		}

		// parse optional flags the command. --due and --priority
		due, _ := cmd.Flags().GetString("due")
		priority, _ := cmd.Flags().GetBool("priority")

		// join all args into a single string to allow task names with spaces
		title := strings.Join(args, " ")

		// if a due date was provided, ensure it uses valid formatting
		if due != "" {
			util.VerifyDate(due)
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
			// exit program and print error
			fmt.Println("Error: Failed to add task", err.Error())
			os.Exit(1)
		}

		// print confirmation
		fmt.Println("Task Created")

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
