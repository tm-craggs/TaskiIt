package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
	"strings"
)

type addFlags struct {
	due      string
	priority bool
}

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
	Use:   "add",
	Short: "add all new task to your to-do list",
	Long:  `Long description goes here`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := task.BackupDB()
		if err != nil {
			return fmt.Errorf("failed to back up database: %w", err)
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no task name provided")
		}

		flags, err := getAddFlags(cmd)
		if err != nil {
			return err
		}

		if flags.due != "" {
			if err := util.VerifyDate(flags.due); err != nil {
				return fmt.Errorf("invalid date format. use YYYY-MM-DD")
			}
		}

		title := strings.Join(args, " ")

		newTask := task.Task{
			Title:    title,
			Due:      flags.due,
			Complete: false,
			Priority: flags.priority,
		}

		if err := task.AddTask(newTask); err != nil {
			return fmt.Errorf("failed to add task: %w", err)
		}

		fmt.Println("Task Created")
		return nil
	},
}

func init() {
	addCmd.Flags().StringP("due", "d", "", "Set all due (e.g. 2025-05-14)")
	addCmd.Flags().BoolP("priority", "p", false, "Mark the task as high priority")

	rootCmd.AddCommand(addCmd)
}
