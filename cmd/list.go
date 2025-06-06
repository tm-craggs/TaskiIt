package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `List all tasks`,

	Run: func(cmd *cobra.Command, args []string) {
		filterPriority, _ := cmd.Flags().GetBool("priority")
		filterComplete, _ := cmd.Flags().GetBool("complete")
		filterNotComplete, _ := cmd.Flags().GetBool("open")
		filterNotPriority, _ := cmd.Flags().GetBool("normal")

		tasks, err := task.GetTasks()
		if err != nil {
			fmt.Println("Failed to get tasks: " + err.Error())
			return
		}

		util.PrintTasks(util.FilterTasks(tasks, filterComplete, filterPriority, filterNotComplete, filterNotPriority))

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("priority", "p", false, "Filter by priority")
	listCmd.Flags().BoolP("complete", "c", false, "Filter by complete")
	listCmd.Flags().BoolP("open", "o", false, "Filter by open task")
	listCmd.Flags().BoolP("normal", "n", false, "Filter by normal priority")
}
