package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "add all new task to your to-do list",
	Long:  `Long description goes here`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !util.ConfirmAction("Confirm Undo?") {
			cmd.SilenceUsage = true
			return fmt.Errorf("aborted by user")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		err := task.RestoreBackup()
		if err != nil {
			return fmt.Errorf("no backup found: %w", err)
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
}
