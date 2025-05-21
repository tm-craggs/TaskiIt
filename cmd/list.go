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

var (
	filterPriority    bool
	filterComplete    bool
	filterNotComplete bool
	filterNotPriority bool
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
		green := p.Color("#00FF00")       // bright green
		red := p.Color("#FF5555")         // light red
		brightBlue := p.Color("#35c5ff ") // bright blue
		orange := p.Color("#FF8000")      // orange
		yellow := p.Color("#FFFF00")      // yellow
		grey := p.Color("#FFFFFF")

		// Helper for bold text
		//bold := func(s string, c termenv.Color) string {
		//	return termenv.String(s).Foreground(c).Bold().String()
		//}

		colour := func(s string, c termenv.Color) string {
			return termenv.String(s).Foreground(c).String()
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"ID", "Title", "Due", "Complete", "Priority"})

		// put tasks through filters
		var filteredTasks []task.Task

		for _, t := range tasks {
			if filterPriority && !t.Priority {
				continue
			}
			if filterComplete && !t.Complete {
				continue
			}
			if filterNotComplete && t.Complete {
				continue
			}
			if filterNotPriority && t.Priority {
				continue
			}
			filteredTasks = append(filteredTasks, t)
		}

		if len(filteredTasks) == 0 {
			fmt.Println("No tasks found")
			return
		}

		for _, t := range filteredTasks {
			var complete, title, due, priority string

			relativeDue := formatDeadline(t.Due)

			if t.Complete {
				// All green if complete
				complete = colour("✔", green)
				title = colour(t.Title, green)
				due = colour(relativeDue, green)
				if t.Priority {
					priority = colour("High", green)
				} else {
					priority = colour("Normal", green)
				}
			} else {
				// Incomplete task

				complete = colour("✘", red)

				// Color title and due based on due date urgency
				switch relativeDue {
				case "Overdue":
					title = colour(t.Title, red)
					due = colour(relativeDue, red)
				case "Today":
					title = colour(t.Title, orange)
					due = colour(relativeDue, orange)
				case "Tomorrow":
					title = colour(t.Title, yellow)
					due = colour(relativeDue, yellow)
				default:
					title = colour(t.Title, grey)
					due = colour(relativeDue, grey)
				}

				// Priority color
				if t.Priority {
					priority = colour("High", brightBlue)
				} else {
					priority = colour("Normal", grey)
				}
			}

			err := table.Append([]string{
				fmt.Sprintf("%d", t.ID),
				title,
				due,
				complete,
				priority,
			})
			if err != nil {
				return
			}
		}

		err = table.Render()
		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&filterPriority, "priority", "p", false, "Filter by priority")
	listCmd.Flags().BoolVarP(&filterComplete, "complete", "c", false, "Filter by complete")
	listCmd.Flags().BoolVarP(&filterNotComplete, "incomplete", "i", false, "Filter by not incomplete")
	listCmd.Flags().BoolVarP(&filterNotPriority, "normal-priority", "n", false, "Filter by normal priority")

}

func formatDeadline(due string) string {
	if due == "" {
		return "None"
	}

	parsedDeadline, err := time.Parse("2006-01-02", due)
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
