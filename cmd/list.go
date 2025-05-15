package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/tcraggs/TidyTask/task"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `List all tasks`,

	Run: func(cmd *cobra.Command, args []string) {

		// get array of all tasks from database
		tasks, err := task.GetTasks()
		if err != nil {
			fmt.Println(err)
			return
		}

		// return if no tasks
		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return
		}

		// define table
		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"ID", "Title", "Deadline", "Complete", "Priority"})

		// iterate through all tasks in array
		for _, t := range tasks {

			// variables defined for entry into table
			var complete, title, deadline, priority string

			// convert date to deadline relative to today
			relativeDeadline := formatDeadline(t.Deadline)

			// if task complete, colour all entries green
			if t.Complete {
				green := color.New(color.FgGreen, color.Bold).SprintFunc()
				complete = green("✔")
				title = green(t.Title)
				deadline = green(relativeDeadline)
				if t.Priority {
					priority = green("High")
				} else {
					priority = green("Normal")
				}
			} else {

				// colour complete red
				complete = color.New(color.FgRed, color.Bold).Sprint("✘")

				// get title
				title = t.Title

				// apply color based on deadline urgency
				switch relativeDeadline {
				case "Overdue":
					deadline = color.New(color.FgRed, color.Bold).Sprint(relativeDeadline)
				case "Today":
					deadline = color.New(color.FgYellow, color.Bold).Sprint(relativeDeadline)
				case "Tomorrow":
					deadline = color.New(color.FgHiYellow, color.Bold).Sprint(relativeDeadline)
				default:
					deadline = relativeDeadline
				}

				// colour priority blue if high
				if t.Priority {
					priority = color.New(color.FgHiBlue, color.Bold).Sprint("High")
				} else {
					priority = "Normal"
				}
			}

			table.Append([]string{
				fmt.Sprintf("%d", t.ID),
				title,
				deadline,
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

func formatDeadline(deadline string) string {
	if deadline == "" {
		return "None"
	}

	parsedDeadline, err := time.Parse("2006-01-02", deadline)
	if err != nil {
		return "Invalid date"
	}

	today := time.Now().Truncate(24 * time.Hour)
	parsedDeadline = parsedDeadline.Truncate(24 * time.Hour)
	difference := parsedDeadline.Sub(today)

	switch days := int(difference.Hours() / 24); {
	case days < 0:
		return "Overdue"
	case days == 0:
		return "Today"
	case days == 1:
		return "Tomorrow"
	case days <= 6:
		return parsedDeadline.Weekday().String()
	default:
		return parsedDeadline.Format("2 Jan 2006")
	}
}
