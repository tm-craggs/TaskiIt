package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"sort"
	"strconv"
	"strings"
)

// create struct that defines the available flags for complete command
type completeFlags struct {
	all      bool
	priority bool
	normal   bool
}

// helper function to parse flags with error handling
func getCompleteFlags(cmd *cobra.Command) (completeFlags, error) {
	var flags completeFlags
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

// completeCmd represents the complete subcommand
var completeCmd = &cobra.Command{
	Use:   "complete [ID...] \n  tidytask complete --all [flags]",
	Short: "Complete tasks from your to-do list",
	Long: `The 'complete' command allows you to mark task as complete in your to-do list.

You can complete tasks in two ways:
1. By specifying one or more task IDs directly
2. Batch completion using the --all flag and optionally applying constraints, such as --priority.
This limits the scope of the complete batch operation to tasks meeting the given criteria. 

You must only use one method. Supplying task IDs together with the --all flag for batch completion causes an error.`,
	Example: `  tidytask complete 1
  > Complete task 1

  tidytask complete 1 2 3
  > Complete tasks 1, 2 and 3

  tidytask complete --all
  > Complete all tasks

  tidytask complete --all --priority
  > Complete all high priority tasks`,

	// main command logic
	RunE: func(cmd *cobra.Command, args []string) error {
		// get flags
		flags, err := getCompleteFlags(cmd)
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

			// create list of task IDs that have been completed
			var completeIDs []int

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

				// complete task
				if err := task.CompleteTask(t.ID); err != nil {
					// add error message to hashmap, with error ID as key value
					failed[t.ID] = err.Error()
				} else {
					// add task ID to complete array if successful
					completeIDs = append(completeIDs, t.ID)
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
				fmt.Println("Failed to complete tasks:")
				for _, id := range keys {
					fmt.Printf("  - %d: %s\n", id, failed[id])
				}
			}

			// throw err if all operations have failed
			if len(completeIDs) == 0 {
				return fmt.Errorf("no tasks completed")
			}

			// define label as tasks, change to task if only one task
			label := "tasks"
			if len(completeIDs) == 1 {
				label = "task"
			}
			// print all successfully completed tasks
			sort.Ints(completeIDs)
			fmt.Printf("Completed %s: %s\n", label, strings.Trim(strings.Replace(fmt.Sprint(completeIDs),
				" ", ", ", -1), "[]"))

			return nil
		}

		// argument given, complete by task IDs

		// backup database
		if err := task.BackupDB(); err != nil {
			return fmt.Errorf("failed to back up database: %w", err)
		}

		// create list of task IDs that have been completed
		var completeIDs []int

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

			// complete task, adding to failed if needed
			if err := task.CompleteTask(id); err != nil {
				failed[strconv.Itoa(id)] = err.Error()
			} else {
				// removal successful, append ID to completed list
				completeIDs = append(completeIDs, id)
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
			fmt.Println("Failed to complete tasks:")
			for _, id := range keys {
				fmt.Printf("  - %s: %s\n", id, failed[id])
			}
		}

		// throw err if all operations have failed
		if len(completeIDs) == 0 {
			return fmt.Errorf("no tasks were completed")
		}

		label := "tasks"
		if len(completeIDs) == 1 {
			label = "task"
		}

		sort.Ints(completeIDs)
		fmt.Printf("Completed %s: %s\n", label, strings.Trim(strings.Replace(fmt.Sprint(completeIDs), " ", ", ", -1), "[]"))

		return nil

	},
}

// command initialisation
func init() {

	// define flags and add subcommand to root

	completeCmd.Flags().BoolP("all", "a", false,
		"Complete all tasks (can be combined with constraints)")

	completeCmd.Flags().BoolP("priority", "p", false,
		"Constrain --all to only complete high priority tasks")

	completeCmd.Flags().BoolP("normal", "n", false,
		"Constrain --all to only complete normal priority tasks")

	rootCmd.AddCommand(completeCmd)
}
