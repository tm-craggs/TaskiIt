package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "add a new task to your to-do list",
	Long:  `Long description goes here`,
	Run: func(cmd *cobra.Command, args []string) {

		tasks, _ := task.SearchTasks("hello", false, true, false)

		keyword := args[0]

		searchID, _ := cmd.Flags().GetBool("id")
		searchTitle, _ := cmd.Flags().GetBool("title")
		searchDue, _ := cmd.Flags().GetBool("due")

		tasks, err := task.SearchTasks(keyword, searchID, searchTitle, searchDue)
		if err != nil {
			fmt.Println("Error searching terms", err)
		}

		filterPriority, _ := cmd.Flags().GetBool("priority")
		filterComplete, _ := cmd.Flags().GetBool("complete")
		filterNotComplete, _ := cmd.Flags().GetBool("open")
		filterNotPriority, _ := cmd.Flags().GetBool("normal")

		util.PrintTasks(util.FilterTasks(tasks, filterComplete, filterPriority, filterNotComplete, filterNotPriority))

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// target flags
	searchCmd.Flags().BoolP("id", "i", false, "Search task by ID")
	searchCmd.Flags().BoolP("title", "t", false, "Search task by title")
	searchCmd.Flags().BoolP("due", "d", false, "Search task by due")

	// filter flags
	searchCmd.Flags().BoolP("complete", "c", false, "Search only complete tasks")
	searchCmd.Flags().BoolP("open", "o", false, "Search only open tasks")
	searchCmd.Flags().BoolP("priority", "p", false, "Search task by priority")
	searchCmd.Flags().BoolP("normal", "n", false, "Search normal priority tasks")

}
