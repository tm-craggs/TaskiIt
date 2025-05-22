package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "add a new task to your to-do list",
	Long:  `Long description goes here`,
	Run: func(cmd *cobra.Command, args []string) {

		taskArr, _ := task.SearchTasks("hello", false, true, false)

		keyword := args[0]

		searchID, _ := cmd.Flags().GetBool("id")
		searchTitle, _ := cmd.Flags().GetBool("title")
		searchDue, _ := cmd.Flags().GetBool("due")

		taskArr, err := task.SearchTasks(keyword, searchID, searchTitle, searchDue)
		if err != nil {
			fmt.Println("Error searching terms", err)
		}

		for _, task := range taskArr {
			fmt.Println(task)
		}

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolP("id", "i", false, "Search task by ID")
	searchCmd.Flags().BoolP("title", "t", false, "Search task by title")
	searchCmd.Flags().BoolP("priority", "p", false, "Search task by priority")
}
