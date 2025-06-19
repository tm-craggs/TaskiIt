package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"sort"
	"strconv"
	"strings"
)

// create struct that defines the available flags for reopen command
type reopenFlags struct {
	all      bool
	priority bool
	normal   bool
}

// helper function to parse flags with error handling
func getReopenFlags(cmd *cobra.Command) (reopenFlags, error) {
	var flags reopenFlags
	var err error

	if flags.all, err = cmd.Flags().GetBool("all"); err != nil {
		return flags, fmt.Errorf("failed to parse --all flag: %w", err)
	}
	if flags.priority, err = cmd.Flags().GetBool("priority"); err != nil {
		return flags, fmt.Errorf("failed to parse --priority flag: %w", err)
	}
	if flags.normal, err = cmd.Flags().GetBool("normal"); err != nil {
		return flags, fmt.Errorf("failed to parse --normal flag: %w", err)
	}

	return flags, nil
}

// reopenCmd represents the reopen subcommand
var reopenCmd = &cobra.Command{
	Use:   "reopen [ID...] \n  tidytask reopen --all [flags]",
	Short: "Reopen tasks from your to-do list",
	Long: `The 'reopen' command allows you to mark complete tasks as open (incomplete) in your to-do list.

You can reopen tasks in two ways:
1. By specifying one or more task IDs directly
2. Batch completion using the --all flag and optionally applying constraints, such as --priority.
This limits the scope of the reopen batch operation to tasks meeting the given criteria. 

You must only use one method. Supplying task IDs together with the --all flag for batch completion causes an error.`,
	Example: `  tidytask reopen 1
  > Reopen task 1

  tidytask reopen 1 2 3
  > Reopen tasks 1, 2 and 3

  tidytask reopen --all
  > Reopen all tasks

  tidytask reopen --all --priority
  > Reopen all high priority tasks`,

	// main command logic
	RunE: func(cmd *cobra.Command, args []string) error {
		// get flags
		flags, err := getReopenFlags(cmd)
		if err != nil {
			return err
		}

		// check flags have not been used with task IDs
		if len(args) > 0 && (flags.all || flags.priority || flags.normal) {
			return fmt.Errorf("cannot use task IDs and batch operation flags together")
		}

		// check for flag conflicts
		if flags.priority && flags.normal {
			return fmt.Errorf("conflicting flags: cannot use --priority and --normal together")
		}

		// filter based removal
		if len(args) == 0 {

			// check for all flag if filter flags have been used
			if !flags.all && (flags.priority || flags.normal) {
				return fmt.Errorf("constraint flags require --all")
			}

			// get tasks
			tasks, err := task.GetTasks()
			if err != nil {
				return fmt.Errorf("failed to retrieve tasks: %w", err)
			}

			// backup database
			if err := task.BackupDB(); err != nil {
				return fmt.Errorf("failed to back up database: %w", err)
			}

			// create list of task IDs that have been reopened
			var reopenIDs []int

			// create hashmap that stores task ID and its error message, if task fails
			failed := make(map[int]string)

			// loop through all tasks
			for _, t := range tasks {

				// check task complies with filters
				if flags.priority && !t.Priority {
					continue
				}
				if flags.normal && t.Priority {
					continue
				}

				// reopen task
				if err := task.ReopenTask(t.ID); err != nil {
					// add error message to hashmap, with error ID as key value
					failed[t.ID] = err.Error()
				} else {
					// add task ID to reopened array if successful
					reopenIDs = append(reopenIDs, t.ID)
				}
			}

			// if there are tasks in failed map, print them to terminal
			if len(failed) > 0 {

				// exact keys and sort
				var keys []int
				for id := range failed {
					keys = append(keys, id)
				}
				sort.Ints(keys)

				// loop through sorted keys array to print failed tasks
				fmt.Println("Failed to reopened tasks:")
				for _, id := range keys {
					fmt.Printf("  - %d: %s\n", id, failed[id])
				}
			}

			// throw err if all operations have failed
			if len(reopenIDs) == 0 {
				return fmt.Errorf("no tasks reopened")
			}

			// define label as tasks, change to task if only one task
			label := "tasks"
			if len(reopenIDs) == 1 {
				label = "task"
			}
			// print all successfully reopened tasks
			sort.Ints(reopenIDs)
			fmt.Printf("Reopened %s: %s\n", label, strings.Trim(strings.Replace(fmt.Sprint(reopenIDs),
				" ", ", ", -1), "[]"))

			return nil
		}

		// argument given, reopen by task IDs

		// backup database
		if err := task.BackupDB(); err != nil {
			return fmt.Errorf("failed to back up database: %w", err)
		}

		// create list of task IDs that have been reopened
		var reopenIDs []int

		// create hashmap that stores task ID as string, and its error message, if task fails
		failed := make(map[string]string)

		// loop through all args input
		for _, arg := range args {

			// parse args into int, add args to failed if needed
			id, err := strconv.Atoi(arg)
			if err != nil {
				failed[arg] = "invalid task ID"
				continue
			}

			// check task exists, add task to failed if task does not exist
			if err := task.CheckTaskExists(id); err != nil {
				failed[strconv.Itoa(id)] = err.Error()
				continue
			}

			// reopen task, adding to failed if needed
			if err := task.ReopenTask(id); err != nil {
				failed[strconv.Itoa(id)] = err.Error()
			} else {
				// removal successful, append ID to reopened list
				reopenIDs = append(reopenIDs, id)
			}
		}

		// if there are tasks in failed map, print them to terminal
		if len(failed) > 0 {

			// exact keys and sort
			var keys []string
			for id := range failed {
				keys = append(keys, id)
			}
			sort.Strings(keys)

			// loop through sorted keys array to print failed tasks
			fmt.Println("Failed to reopen tasks:")
			for _, id := range keys {
				fmt.Printf("  - %s: %s\n", id, failed[id])
			}
		}

		// throw err if all operations have failed
		if len(reopenIDs) == 0 {
			return fmt.Errorf("no tasks were reopened")
		}

		label := "tasks"
		if len(reopenIDs) == 1 {
			label = "task"
		}

		sort.Ints(reopenIDs)
		fmt.Printf("Reopened %s: %s\n", label, strings.Trim(strings.Replace(fmt.Sprint(reopenIDs), " ", ", ", -1), "[]"))

		return nil

	},
}

// command initialisation
func init() {

	// define flags and add subcommand to root

	reopenCmd.Flags().BoolP("all", "a", false,
		"Reopen all tasks (can be combined with constraints)")

	reopenCmd.Flags().BoolP("priority", "p", false,
		"Constrain --all to only reopen high priority tasks")

	reopenCmd.Flags().BoolP("normal", "n", false,
		"Constrain --all to only reopen normal priority tasks")

	rootCmd.AddCommand(reopenCmd)
}
