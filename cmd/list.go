package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/muesli/termenv"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
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

		// Setup termenv profile for colors
		p := termenv.ColorProfile()

		// Define colors (using RGB for more precise control)
		green := p.Color("#00FF00")      // bright green
		red := p.Color("#FF5555")        // light red
		white := p.Color("#FFFFFF")      // white
		brightBlue := p.Color("#569CD6") // bright blue
		orange := p.Color("#FF8000")     // orange
		yellow := p.Color("#FFFF00")     // yellow

		// Helper for bold text
		bold := func(s string, c termenv.Color) string {
			return termenv.String(s).Foreground(c).Bold().String()
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"ID", "Title", "Due", "Complete", "Priority"})

		for _, t := range tasks {
			var complete, title, due, priority string

			relativeDue := formatDeadline(t.Deadline)

			if t.Complete {
				// All green if complete
				complete = bold("✔", green)
				title = bold(t.Title, green)
				due = bold(relativeDue, green)
				if t.Priority {
					priority = bold("High", green)
				} else {
					priority = bold("Normal", green)
				}
			} else {
				// Incomplete task

				complete = bold("✘", red)

				// Color title and due based on due date urgency
				switch relativeDue {
				case "Overdue":
					title = bold(t.Title, red)
					due = bold(relativeDue, red)
				case "Today":
					title = bold(t.Title, orange)
					due = bold(relativeDue, orange)
				case "Tomorrow":
					title = bold(t.Title, yellow)
					due = bold(relativeDue, yellow)
				default:
					title = bold(t.Title, white)
					due = bold(relativeDue, white)
				}

				// Priority color
				if t.Priority {
					priority = bold("High", brightBlue)
				} else {
					priority = bold("Normal", white)
				}
			}

			table.Append([]string{
				fmt.Sprintf("%d", t.ID),
				title,
				due,
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
