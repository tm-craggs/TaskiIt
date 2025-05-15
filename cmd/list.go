package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	//"github.com/spf13/pflag"
	"github.com/tcraggs/TidyTask/task"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `List all tasks`,

	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := task.GetTasks()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(tasks)

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
