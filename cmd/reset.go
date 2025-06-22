package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tm-craggs/tidytask/task"
	"github.com/tm-craggs/tidytask/util"
)

var resetCmd = &cobra.Command{
	Use:                   "reset",
	DisableFlagsInUseLine: true,
	Short:                 "Permanently delete all task data and backups",
	Long: `The 'reset' command hard resets TidyTask by deleting the file which stores your tasks, and the backup file
if they exist. 

This is useful to recover from corrupted files, or to restart task ID numbering at 1

This command cannot be reversed once run.`,

	// confirm action before running command
	PreRunE: func(cmd *cobra.Command, args []string) error {

		// check args
		if len(args) > 0 {
			return fmt.Errorf("unexpected arguments: %v; use --help for usage information", args)
		}

		if !util.ConfirmAction("Confirm hard reset?") {
			return fmt.Errorf("aborted by user")
		}
		return nil
	},

	// command logic
	RunE: func(cmd *cobra.Command, args []string) error {

		// attempt to restore backup, throw error if no backup file found
		err := task.HardReset()
		if err != nil {
			return fmt.Errorf("failed to hard reset: %w", err)
		}

		// exit
		return nil

	},
}

// command initialisation
func init() {

	// add subcommand to root
	rootCmd.AddCommand(resetCmd)
}
