package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tm-craggs/tidytask/task"
	"github.com/tm-craggs/tidytask/util"
	"sort"
	"strconv"
	"strings"
)

// create struct that defines the available flags for remove command
type removeFlags struct {
	all      bool
	complete bool
	priority bool
	open     bool
	normal   bool
}

// helper function to parse flags with error handling
func getRemoveFlags(cmd *cobra.Command) (removeFlags, error) {
	var flags removeFlags
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
	if flags.open, err = cmd.Flags().GetBool("open"); err != nil {
		return flags, fmt.Errorf("failed to parse --open flag: %w", err)
	}
	if flags.complete, err = cmd.Flags().GetBool("complete"); err != nil {
		return flags, fmt.Errorf("failed to parse --complete flag: %w", err)
	}

	return flags, nil
}

// removeCmd represents the remove subcommand
var removeCmd = &cobra.Command{
	Use:   "remove [ID...] \n  tidytask remove --all [flags]",
	Short: "Remove tasks from your to-do list",
	Long: `The 'remove' command allows you to delete tasks from your to-do list.

You can remove tasks in two ways:
1. By specifying one or more task IDs directly
2. Batch removal using the --all flag and optionally applying constraints, such as --priority.
This limits the scope of the removal batch operation to tasks meeting the given criteria. 

You must only use one method. Supplying task IDs together with the --all flag for batch removal causes an error.

When combining constraints, such as --priority and --complete, it will only remove tasks that meet all conditions.`,
	Example: `  tidytask remove 1
  > Remove task 1

  tidytask remove 1 2 3
  > Remove tasks 1, 2 and 3

  tidytask remove --all
  > Remove all tasks

  tidytask remove --all --priority
  > Remove all high priority tasks

  tidytask remove --all --priority --complete
  > Remove all tasks that are high priority AND complete`,

	// main command logic
	RunE: func(cmd *cobra.Command, args []string) error {

		// get flags
		flags, err := getRemoveFlags(cmd)
		if err != nil {
			return err
		}

		// disallow no input
		if len(args) == 0 && !flags.all {
			return fmt.Errorf("no arguments provided; task ID or --all flag required")
		}

		// disallow mixed usage
		if len(args) > 0 && (flags.all || flags.priority || flags.normal) {
			return fmt.Errorf("cannot use task IDs and batch operation flags together")
		}

		// check flags have not been used with task IDs
		if len(args) > 0 && (flags.all || flags.priority || flags.normal || flags.complete || flags.open) {
			return fmt.Errorf("cannot use task IDs and batch operation flags together")
		}

		// check for flag conflicts
		if flags.priority && flags.normal {
			return fmt.Errorf("conflicting flags: cannot use --priority and --normal together")
		}

		if flags.complete && flags.open {
			return fmt.Errorf("conflicting flags: cannot use --complete and --open together")
		}

		// filter based removal
		if len(args) == 0 {

			// check for all flag if filter flags have been used
			if !flags.all && (flags.priority || flags.normal || flags.complete || flags.open) {
				return fmt.Errorf("constraint flags require --all")
			}

			// get tasks
			tasks, err := task.GetTasks()
			if err != nil {
				return fmt.Errorf("failed to retrieve tasks: %w", err)
			}

			// backup database
			if err := task.BackupDB(); err != nil {
				fmt.Printf("Warning: failed to back up database: %v", err)
			}

			// prompt for confirmation
			if !util.ConfirmAction("Confirm Removal?") {
				cmd.SilenceUsage = true
				return fmt.Errorf("aborted by user")
			}

			// create list of task IDs that have been removed
			var removedIDs []int

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
				if flags.complete && !t.Complete {
					continue
				}
				if flags.open && t.Complete {
					continue
				}

				// remove task
				if err := task.RemoveTask(t.ID); err != nil {
					// add error message to hashmap, with error ID as key value
					failed[t.ID] = err.Error()
				} else {
					// add task ID to removed array if successful
					removedIDs = append(removedIDs, t.ID)
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
				fmt.Println("Failed to remove tasks:")
				for _, id := range keys {
					fmt.Printf("  - %d: %s\n", id, failed[id])
				}
			}

			// throw err if all operations have failed
			if len(removedIDs) == 0 {
				return fmt.Errorf("no tasks removed")
			}

			// define label as tasks, change to task if only one task
			label := "tasks"
			if len(removedIDs) == 1 {
				label = "task"
			}
			// print all successfully removed tasks
			sort.Ints(removedIDs)
			fmt.Printf("Removed %s: %s\n", label, strings.Trim(strings.Replace(fmt.Sprint(removedIDs),
				" ", ", ", -1), "[]"))

			return nil
		}

		// argument given, remove by task IDs

		// backup database
		if err := task.BackupDB(); err != nil {
			return fmt.Errorf("failed to back up database: %w", err)
		}

		// prompt for confirmation
		if !util.ConfirmAction("Confirm Removal?") {
			cmd.SilenceUsage = true
			return fmt.Errorf("aborted by user")
		}

		// create list of task IDs that have been removed
		var removedIDs []int

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

			// remove task, adding to failed if needed
			if err := task.RemoveTask(id); err != nil {
				failed[strconv.Itoa(id)] = err.Error()
			} else {
				// removal successful, append ID to removed list
				removedIDs = append(removedIDs, id)
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
			fmt.Println("Failed to remove tasks:")
			for _, id := range keys {
				fmt.Printf("  - %s: %s\n", id, failed[id])
			}
		}

		// throw err if all operations have failed
		if len(removedIDs) == 0 {
			return fmt.Errorf("no tasks were removed")
		}

		label := "tasks"
		if len(removedIDs) == 1 {
			label = "task"
		}

		sort.Ints(removedIDs)
		fmt.Printf("Removed %s: %s\n", label, strings.Trim(strings.Replace(fmt.Sprint(removedIDs), " ", ", ", -1), "[]"))

		return nil

	},
}

// command initialisation
func init() {

	// define flags and add subcommand to root

	removeCmd.Flags().BoolP("all", "a", false,
		"Remove all tasks (can be combined with constraints)")

	removeCmd.Flags().BoolP("priority", "p", false,
		"Constrain --all to only remove high priority tasks")

	removeCmd.Flags().BoolP("normal", "n", false,
		"Constrain --all to only remove normal priority tasks")

	removeCmd.Flags().BoolP("complete", "c", false,
		"Constrain --all to only remove complete tasks")

	removeCmd.Flags().BoolP("open", "o", false,
		"Constrain --all to only remove open (incomplete) tasks")

	rootCmd.AddCommand(removeCmd)
}
