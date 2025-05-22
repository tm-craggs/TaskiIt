package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "add a new task to your to-do list",
	Long:  `Long description goes here`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !util.ConfirmAction("Confirm Undo?") {
			cmd.SilenceUsage = true
			return fmt.Errorf("aborted by user")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		err := task.RestoreBackup()
		if err != nil {
			fmt.Println("No backup found")
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
}
