package util

import (
	"fmt"
	"github.com/muesli/termenv"
	"github.com/olekukonko/tablewriter"
	"github.com/tcraggs/TidyTask/task"
	"os"
	"time"
)

func PrintTasks(tasks []task.Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	p := termenv.ColorProfile()

	green := p.Color("#00FF00")
	red := p.Color("#FF5555")
	brightBlue := p.Color("#35c5ff")
	orange := p.Color("#FF8000")
	yellow := p.Color("#FFFF00")
	grey := p.Color("#FFFFFF")

	colour := func(s string, c termenv.Color) string {
		return termenv.String(s).Foreground(c).String()
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Title", "Due", "Complete", "Priority"})

	for _, t := range tasks {
		var complete, title, due, priority string

		if t.Complete {
			complete = colour("✔", green)
			title = colour(t.Title, green)

			if t.Due == "" {
				due = colour("Met: No due", green)
			} else {
				dueDate, err1 := time.Parse("2006-01-02", t.Due)
				var completeDate time.Time
				var err2 error

				if t.CompleteDate.Valid {
					completeDate, err2 = time.Parse("2006-01-02", t.CompleteDate.String)
				} else {
					err2 = fmt.Errorf("no complete date")
				}

				if err1 == nil && err2 == nil {
					diff := int(completeDate.Truncate(24*time.Hour).Sub(dueDate.Truncate(24*time.Hour)).Hours() / 24)
					diffText := dateDiff(dueDate, completeDate)

					switch {
					case diff == 0:
						due = "Met: On Time"
					case diff < 0:
						due = fmt.Sprintf("Met: %s early", diffText)
					default:
						due = fmt.Sprintf("Missed: %s late", diffText)
					}
					due = colour(due, green)
				} else {
					due = colour(t.Due, green)
				}
			}

			if t.Priority {
				priority = colour("High", green)
			} else {
				priority = colour("Normal", green)
			}
		} else {
			relativeDue := formatDeadline(t.Due)
			complete = colour("✘", red)

			// this code is horrible and needs to be nuked
			// TODO: Append overdue with time in the dateDiff function

			if relativeDue == "Overdue" {
				dueDate, err := time.Parse("2006-01-02", t.Due)
				if err == nil {
					diffText := dateDiff(dueDate, time.Now())
					due = colour(fmt.Sprintf("Overdue: %s", diffText), red)
				} else {
					due = colour(relativeDue, red)
				}
				title = colour(t.Title, red)
			} else {
				switch relativeDue {
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
			}

			if t.Priority {
				priority = colour("High", brightBlue)
			} else {
				priority = colour("Normal", grey)
			}
		}

		if err := table.Append([]string{
			fmt.Sprintf("%d", t.ID),
			title,
			due,
			complete,
			priority,
		}); err != nil {
			return
		}
	}

	_ = table.Render()
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
	days := int(parsedDeadline.Sub(today).Hours() / 24)

	switch {
	case days < 0:
		return "Overdue"
	case days == 0:
		return "Today"
	case days == 1:
		return "Tomorrow"
	case days <= 6:
		return parsedDeadline.Weekday().String()
	default:
		return parsedDeadline.Format("2006-01-02")
	}
}
