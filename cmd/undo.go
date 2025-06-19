package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tm-craggs/tidytask/task"
	"github.com/tm-craggs/tidytask/util"
)

var undoCmd = &cobra.Command{
	Use:                   "undo",
	DisableFlagsInUseLine: true,
	Short:                 "Undo the previous action",
	Long: `The 'undo' command restores the most recent backup of your task list.

A backup is automatically created before any command that modifies tasks, restoring the backup will reverse that change. 

Only one backup is stored at a time - it can only restore the most recent change

After running the undo command, it becomes unavailable until a new change is made.`,

	// confirm action before running command
	PreRunE: func(cmd *cobra.Command, args []string) error {

		// check args
		if len(args) > 0 {
			return fmt.Errorf("unexpected arguments: %v; use --help for usage information", args)
		}

		if !util.ConfirmAction("Confirm Undo?") {
			return fmt.Errorf("aborted by user")
		}
		return nil
	},

	// command logic
	RunE: func(cmd *cobra.Command, args []string) error {

		// attempt to restore backup, throw error if no backup file found
		err := task.RestoreBackup()
		if err != nil {
			return fmt.Errorf("no backup found: %w", err)
		}

		// exit
		return nil

	},
}

// command initialisation
func init() {

	// add subcommand to root
	rootCmd.AddCommand(undoCmd)
}
