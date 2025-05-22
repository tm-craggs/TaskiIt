package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"strconv"
)

var (
	editDue   string
	editTitle string
)

// addCmd represents the add command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "add a new task to your to-do list",
	Long:  `Long description goes here`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := task.BackupDB()
		if err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please enter a task name")
			return
		}

		// task title is taken in as args
		id, _ := strconv.Atoi(args[0])

		if cmd.Flags().Changed("title") {
			if err := task.SetTitle(id, editTitle); err != nil {
				fmt.Println(err)
			}
		}

		if cmd.Flags().Changed("priority") {
			if err := task.TogglePriority(id); err != nil {
				fmt.Println(err)
			}
		}

		if cmd.Flags().Changed("due") {
			if err := task.SetDue(id, editDue); err != nil {
				fmt.Println(err)
			}
		}

		fmt.Println("Task updated successfully")

	},
}

func init() {
	editCmd.Flags().StringVarP(&editDue, "due", "d", "", "Set a due date (e.g. 2025-05-14)")
	editCmd.Flags().BoolVarP(priority, "priority", "p", false, "Mark the task as high priority")
	editCmd.Flags().StringVarP(&editTitle, "title", "t", "", "Set a task title")

	rootCmd.AddCommand(editCmd)

}
