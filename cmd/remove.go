package cmd

// TODO: Rework to use filters, pass filtered tasks into a function that removes an array of Tasks.

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
	"strconv"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a task",
	Long:  `Long goes here`,

	// confirmation for delete as pre-run
	PreRunE: func(cmd *cobra.Command, args []string) error {

		removeAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}

		removeComplete, err := cmd.Flags().GetBool("complete")
		if err != nil {
			return err
		}

		if removeAll && removeComplete {
			return fmt.Errorf("conflicting flags: --complete cannot be used with --all")
		}

		var confirmationText string

		switch {
		case removeAll:
			confirmationText = "Remove all tasks?"
		case removeComplete:
			confirmationText = "Remove all completed tasks?"
		default:
			if len(args) == 0 {
				return fmt.Errorf("task ID required")
			}
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

	// core logic for `complete` command
	RunE: func(cmd *cobra.Command, args []string) error {

		// parse flags
		removeAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			return fmt.Errorf("failed to get flag options: %w", err)
		}

		removeComplete, err := cmd.Flags().GetBool("complete")
		if err != nil {
			return fmt.Errorf("failed to get flag options: %w", err)
		}

		// code for --all flag
		if removeAll {

			if len(args) > 0 {
				return fmt.Errorf("--all flag cannot be used with task IDs")
			}

			if err := task.RemoveAllTasks(); err != nil {
				return fmt.Errorf("failed to remove all tasks: %w", err)
			}

			// print confirmation and return
			fmt.Println("All tasks removed")
			return nil
		}

		// code for --complete flags
		if removeComplete {

			if len(args) > 0 {
				return fmt.Errorf("--complete cannot be used with task IDs")
			}

			if err := task.RemoveAllCompleteTasks(); err != nil {
				return fmt.Errorf("failed to remove all complete tasks: %w", err)
			}

			// print confirmation and return
			fmt.Println("All complete tasks removed")
			return nil

		}

		// single task removal for single logic

		// check input has been given
		if len(args) == 0 {
			return fmt.Errorf("task ID required")
		}

		// convert args input to int
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("failed to parse task ID: %w", err)
		}

		// check if input ID exists in DB
		if err := task.CheckTaskExists(id); err != nil {
			return fmt.Errorf("task does not exist: %w", err)
		}

		// call complete task
		if err := task.RemoveTask(id); err != nil {
			return fmt.Errorf("failed to remove task: %w", err)
		}

		// print confirmation
		fmt.Println("Task removed")
		return nil
	},
}

func init() {
	removeCmd.Flags().BoolP("all", "a", false, "Remove all tasks")
	removeCmd.Flags().BoolP("complete", "c", false, "Remove all complete tasks")
	rootCmd.AddCommand(removeCmd)
}
