package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
	"strconv"
	"strings"
)

type removeFlags struct {
	all      bool
	complete bool
	priority bool
	open     bool
	normal   bool
}

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

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove tasks",
	Long:  `Remove one or more tasks by ID, or use filters with --all to bulk-remove.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := getRemoveFlags(cmd)
		if err != nil {
			return err
		}

		// Case 1: No args — filter-based or full removal
		if len(args) == 0 {
			if !flags.all && (flags.priority || flags.normal || flags.complete || flags.open) {
				return fmt.Errorf("filter flags require --all")
			}

			if !flags.all {
				return fmt.Errorf("task ID required")
			}

			tasks, err := task.GetTasks()
			if err != nil {
				return fmt.Errorf("failed to retrieve tasks: %w", err)
			}

			if err := task.BackupDB(); err != nil {
				return fmt.Errorf("failed to back up database: %w", err)
			}
			if !util.ConfirmAction("Confirm Removal?") {
				cmd.SilenceUsage = true
				return fmt.Errorf("aborted by user")
			}

			var removedIDs []int
			failed := make(map[int]string)

			for _, t := range tasks {
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
				if err := task.RemoveTask(t.ID); err != nil {
					failed[t.ID] = err.Error()
				} else {
					removedIDs = append(removedIDs, t.ID)
				}
			}

			if len(removedIDs) == 0 {
				return fmt.Errorf("no tasks matched specified filters")
			}

			label := "tasks"
			if len(removedIDs) == 1 {
				label = "task"
			}
			fmt.Printf("Removed %s: %s\n", label, strings.Trim(strings.Replace(fmt.Sprint(removedIDs), " ", ", ", -1), "[]"))

			if len(failed) > 0 {
				fmt.Println("failed to remove tasks:")
				for id, reason := range failed {
					fmt.Printf("  - %d: %s\n", id, reason)
				}
			}

			return nil
		}

		// Case 2: Arguments given — remove by task IDs
		if err := task.BackupDB(); err != nil {
			return fmt.Errorf("failed to back up database: %w", err)
		}
		if !util.ConfirmAction("Confirm Removal?") {
			cmd.SilenceUsage = true
			return fmt.Errorf("aborted by user")
		}

		var removedIDs []int
		failed := make(map[int]string)

		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("invalid task ID '%s': %v\n", arg, err)
				continue
			}

			if err := task.CheckTaskExists(id); err != nil {
				failed[id] = err.Error()
				continue
			}

			if err := task.RemoveTask(id); err != nil {
				failed[id] = err.Error()
			} else {
				removedIDs = append(removedIDs, id)
			}
		}

		if len(removedIDs) == 0 {
			return fmt.Errorf("no tasks were removed")
		}

		label := "tasks"
		if len(removedIDs) == 1 {
			label = "task"
		}
		fmt.Printf("Removed %s: %s\n", label, strings.Trim(strings.Replace(fmt.Sprint(removedIDs), " ", ", ", -1), "[]"))

		if len(failed) > 0 {
			fmt.Println("failed to remove tasks:")
			for id, reason := range failed {
				fmt.Printf("  - %d: %s\n", id, reason)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().Bool("all", false, "Apply to all tasks (required for filters or to remove all)")
	removeCmd.Flags().Bool("priority", false, "Only remove priority tasks")
	removeCmd.Flags().Bool("normal", false, "Only remove non-priority tasks")
	removeCmd.Flags().Bool("complete", false, "Only remove completed tasks")
	removeCmd.Flags().Bool("open", false, "Only remove open tasks")
}
