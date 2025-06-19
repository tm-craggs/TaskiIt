package util

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/muesli/termenv"         // for colour styling in terminal output
	"github.com/olekukonko/tablewriter" // for rendering tables in terminal
	"github.com/tcraggs/tidytask/task"
)

var (
	// set up terminal colour profile and predefined colour values
	p          = termenv.ColorProfile()
	green      = p.Color("#00FF00")
	red        = p.Color("#FF5555")
	brightBlue = p.Color("#35c5ff")
	orange     = p.Color("#FF8000")
	yellow     = p.Color("#FFFF00")
	grey       = p.Color("#FFFFFF")

	// helper function to apply colour to strings
	colorise = func(s string, c termenv.Color) string {
		return termenv.String(s).Foreground(c).String()
	}

	// helper function to truncate time to midnight for date only comparison
	truncateTime = func(t time.Time) time.Time {
		return t.Truncate(24 * time.Hour)
	}
)

// PrintTasks takes a slice of Task structs and displays them as a colour coded table in the terminal
func PrintTasks(tasks []task.Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	// create table and set up table headers
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "title", "due", "complete", "priority"})

	// iterate through all tasks in the slice
	for _, t := range tasks {

		// create empty strings for each column
		var title, due, complete, priority string

		if t.Complete {
			// if task is complete, format task as complete, and colour priority green
			// due will be formatted on whether it was on time, late, or early.
			complete, title, due = formatCompletedTask(t)
			priority = formatPriority(t.Priority, green, green)
		} else {
			// if task is incomplete, format as incomplete
			complete, title, due = formatIncompleteTask(t)
			// incomplete tasks use brightBlue for priority, grey for normal
			priority = formatPriority(t.Priority, brightBlue, grey)
		}

		// append the formatted task data as a row in the table
		if err := table.Append([]string{
			fmt.Sprintf("%d", t.ID), // task ID as string
			title,                   // coloured task title
			due,                     // stylised due date
			complete,                // tick or cross with colour
			priority,                // high or normal with colour
		}); err != nil {
			// if appending fails, log and move to next task
			log.Printf("Error: Failed to append task ID %d to table: %v", t.ID, err)
			continue
		}
	}

	// render final table
	err := table.Render()
	if err != nil {
		log.Fatal(err)
	}
}

// formatCompletedTask returns styled fields for a completed task
func formatCompletedTask(t task.Task) (string, string, string) {

	// colour green tick for complete field
	complete := colorise("✔", green)

	// colour title field green
	title := colorise(t.Title, green)

	// if due date empty, display "Met: No due"
	if t.Due == "" {
		return complete, title, colorise("Met: No due", green)
	}

	// get current date
	dueDate, err1 := time.Parse("2006-01-02", t.Due)

	// get complete date
	completeDate, err2 := time.Parse("2006-01-02", t.CompleteDate.String)

	// if either parsing fails, display raw due date without being relative.
	if err1 != nil || err2 != nil {
		return complete, title, colorise(t.Due, green)
	}

	// calculate the difference in full days between due date and completion date
	diff := int(truncateTime(completeDate).Sub(truncateTime(dueDate)).Hours() / 24)

	// get a human-readable difference between the dates
	diffText := dateDiff(dueDate, completeDate)

	// format the due based on difference between due and submitted
	var due string
	switch {
	case diff == 0:
		due = "Met: On Time"
	case diff < 0:
		due = fmt.Sprintf("Met: %s early", diffText)
	default:
		due = fmt.Sprintf("Missed: %s late", diffText)
	}

	// return 3 formatted fields
	return complete, title, colorise(due, green)
}

func formatIncompleteTask(t task.Task) (string, string, string) {

	// colour green tick for complete field
	complete := colorise("✘", red)

	// get a relative due date
	relativeDue := formatDeadline(t.Due)

	// if task overdue, colour complete, title and due as red
	if strings.HasPrefix(relativeDue, "Overdue") {
		return complete, colorise(t.Title, red), colorise(relativeDue, red)
	}

	// colour due date based on closeness of due
	switch relativeDue {
	case "Today":
		return complete, colorise(t.Title, orange), colorise(relativeDue, orange)
	case "Tomorrow":
		return complete, colorise(t.Title, yellow), colorise(relativeDue, yellow)
	default:
		return complete, colorise(t.Title, grey), colorise(relativeDue, grey)
	}
}

// formatPriority returns a styled string representing the task's priority level.
func formatPriority(isHigh bool, highColor, normalColor termenv.Color) string {
	// if the task is high priority, return "High" styled with the highColor.
	if isHigh {
		return colorise("High", highColor)
	}

	// if the task is not high priority, return "Normal" styled with the normalColor.
	return colorise("Normal", normalColor)
}

// formatDeadline formats a due date string into a human-readable status.
// it returns "None", "Today", "Tomorrow", a weekday name, or an ISO date.
// if the date is past, it returns an "Overdue" label with how long it's overdue.
func formatDeadline(due string) string {

	// if there is no due date, return "None"
	if due == "" {
		return "None"
	}

	// parse due date using layout YYYY-MM-DD
	parsedDue, err := time.Parse("2006-01-02", due)

	// if parsing fails, return "Invalid date"
	if err != nil {
		return "Invalid date"
	}

	// get today's date, truncated to remove time
	today := time.Now().Truncate(24 * time.Hour)

	// truncate parsedDue to remove time
	parsedDue = parsedDue.Truncate(24 * time.Hour)

	// calculate the difference in days between today's date and the due date
	days := int(parsedDue.Sub(today).Hours() / 24)

	// show overdue, today, or tomorrow, or day if the week when task is within a week.
	// otherwise, show raw date
	switch {
	case days < 0:
		diffText := dateDiff(parsedDue, today)
		return fmt.Sprintf("Overdue: %s", diffText)
	case days == 0:
		return "Today"
	case days == 1:
		return "Tomorrow"
	case days <= 6:
		return parsedDue.Weekday().String()
	default:
		return parsedDue.Format("2006-01-02")
	}
}
