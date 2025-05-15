package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	//"github.com/spf13/pflag"
	"github.com/tcraggs/TidyTask/task"
	"strconv"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a task",
	Long:  `Long goes here`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println("Please specify the task ID to remove")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Task ID not found")
		}

		if err := task.DeleteTask(id); err != nil {
			fmt.Println("Failed to remove task:", err.Error())
			return
		}

		fmt.Println("Task removed successfully")

	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
