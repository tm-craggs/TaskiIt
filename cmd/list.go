package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	//"github.com/spf13/pflag"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
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

		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"ID", "Title", "Deadline", "Complete", "Priority"})

		for _, t := range tasks {
			var complete string
			if t.Complete {
				complete = color.New(color.FgGreen, color.Bold).Sprint("✔")
			} else {
				complete = color.New(color.FgRed, color.Bold).Sprint("✘")
			}

			priority := "Normal"
			if t.Priority {
				priority = "High"
			}

			table.Append([]string{
				fmt.Sprintf("%d", t.ID),
				t.Title,
				t.Deadline,
				complete,
				priority,
			})

		}

		table.Render()

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
