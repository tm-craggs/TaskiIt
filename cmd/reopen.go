package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"strconv"
	"strings"
)

type reopenFlags struct {
	all      bool
	priority bool
	normal   bool
}

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

var reopenCmd = &cobra.Command{
	Use:   "reopen",
	Short: "Reopen tasks",
	Long:  `Mark completed tasks as incomplete.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := getReopenFlags(cmd)
		if err != nil {
			return err
		}

		// Case 1: No args — filter-based or full removal
		if len(args) == 0 {
			if !flags.all && (flags.priority || flags.normal) {
				return fmt.Errorf("filter flags require --all")
			}

			if !flags.all {
				return fmt.Errorf("task ID required")
			}

			tasks, err := task.GetTasks()
			if err != nil {
				return fmt.Errorf("failed to retrieve tasks: %w", err)
			}

			var reopenIDs []int
			failed := make(map[int]string)

			for _, t := range tasks {
				if flags.priority && !t.Priority {
					continue
				}
				if flags.normal && t.Priority {
					continue
				}

				if err := task.ReopenTask(t.ID); err != nil {
					failed[t.ID] = err.Error()
				} else {
					reopenIDs = append(reopenIDs, t.ID)
				}
			}

			if len(reopenIDs) == 0 {
				return fmt.Errorf("no tasks matched specified filters")
			}

			label := "tasks"
			if len(reopenIDs) == 1 {
				label = "task"
			}
			fmt.Printf("Reopened %s: %s\n", label, strings.Trim(strings.Replace(fmt.Sprint(reopenIDs), " ", ", ", -1), "[]"))

			if len(failed) > 0 {
				fmt.Println("failed to reopen tasks:")
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

		var reopenIDs []int
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

			if err := task.ReopenTask(id); err != nil {
				failed[id] = err.Error()
			} else {
				reopenIDs = append(reopenIDs, id)
			}
		}

		if len(reopenIDs) == 0 {
			return fmt.Errorf("no tasks were reopened")
		}

		label := "tasks"
		if len(reopenIDs) == 1 {
			label = "task"
		}
		fmt.Printf("Reopened %s: %s\n", label, strings.Trim(strings.Replace(fmt.Sprint(reopenIDs), " ", ", ", -1), "[]"))

		if len(failed) > 0 {
			fmt.Println("failed to reopen tasks:")
			for id, reason := range failed {
				fmt.Printf("  - %d: %s\n", id, reason)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(reopenCmd)

	reopenCmd.Flags().Bool("all", false, "Apply to all tasks (required for filters or to remove all)")
	reopenCmd.Flags().Bool("priority", false, "Only remove priority tasks")
	reopenCmd.Flags().Bool("normal", false, "Only remove non-priority tasks")

}
