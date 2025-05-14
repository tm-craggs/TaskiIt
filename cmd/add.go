package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	message  string
	deadline string
	priority bool
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new task to your to-do list",
	Long:  `Long description goes here`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please enter a task name")
		}
		title := args[0]
		fmt.Println(title)

		if message != "" {
			fmt.Println(message)
		}
		if deadline != "" {
			fmt.Println(deadline)
		}
		if priority {
			fmt.Println(priority)
		}
	},
}

func init() {
	addCmd.Flags().StringVarP(&message, "message", "m", "", "Message to send")
	addCmd.Flags().StringVarP(&deadline, "deadline", "d", "", "Set a deadline (e.g. 2025-05-14)")
	addCmd.Flags().BoolVarP(&priority, "priority", "p", false, "Mark the task as high priority")

	rootCmd.AddCommand(addCmd)

}
