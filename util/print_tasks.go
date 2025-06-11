package util

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/muesli/termenv"
	"github.com/olekukonko/tablewriter"
	"github.com/tcraggs/TidyTask/task"
)

var (
	p            = termenv.ColorProfile()
	green        = p.Color("#00FF00")
	red          = p.Color("#FF5555")
	brightBlue   = p.Color("#35c5ff")
	orange       = p.Color("#FF8000")
	yellow       = p.Color("#FFFF00")
	grey         = p.Color("#FFFFFF")
	colorise     = func(s string, c termenv.Color) string { return termenv.String(s).Foreground(c).String() }
	truncateTime = func(t time.Time) time.Time { return t.Truncate(24 * time.Hour) }
)

func PrintTasks(tasks []task.Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "title", "due", "complete", "priority"})

	for _, t := range tasks {
		var title, due, complete, priority string

		if t.Complete {
			complete, title, due = formatCompletedTask(t)
			priority = formatPriority(t.Priority, green, green)
		} else {
			complete, title, due = formatIncompleteTask(t)
			priority = formatPriority(t.Priority, brightBlue, grey)
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

	err := table.Render()
	if err != nil {
		log.Fatal(err)
	}
}

func formatCompletedTask(t task.Task) (string, string, string) {
	complete := colorise("✔", green)
	title := colorise(t.Title, green)

	if t.Due == "" {
		return complete, title, colorise("Met: No due", green)
	}

	dueDate, err1 := time.Parse("2006-01-02", t.Due)
	completeDate, err2 := time.Parse("2006-01-02", t.CompleteDate.String)

	if err1 != nil || err2 != nil {
		return complete, title, colorise(t.Due, green)
	}

	diff := int(truncateTime(completeDate).Sub(truncateTime(dueDate)).Hours() / 24)
	diffText := dateDiff(dueDate, completeDate)

	var due string
	switch {
	case diff == 0:
		due = "Met: On Time"
	case diff < 0:
		due = fmt.Sprintf("Met: %s early", diffText)
	default:
		due = fmt.Sprintf("Missed: %s late", diffText)
	}

	return complete, title, colorise(due, green)
}

func formatIncompleteTask(t task.Task) (string, string, string) {
	complete := colorise("✘", red)
	relativeDue := formatDeadline(t.Due)

	if strings.HasPrefix(relativeDue, "Overdue") {
		return complete, colorise(t.Title, red), colorise(relativeDue, red)
	}

	switch relativeDue {
	case "Today":
		return complete, colorise(t.Title, orange), colorise(relativeDue, orange)
	case "Tomorrow":
		return complete, colorise(t.Title, yellow), colorise(relativeDue, yellow)
	default:
		return complete, colorise(t.Title, grey), colorise(relativeDue, grey)
	}
}

func formatPriority(isHigh bool, highColor, normalColor termenv.Color) string {
	if isHigh {
		return colorise("High", highColor)
	}
	return colorise("Normal", normalColor)
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
		diffText := dateDiff(parsedDeadline, today)
		return fmt.Sprintf("Overdue: %s", diffText)
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
