package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"strconv"
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
	Long:  `Remove a task, by ID or using filters`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := task.BackupDB(); err != nil {
			return fmt.Errorf("failed to back up database: %w", err)
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := getRemoveFlags(cmd) // fixed function name
		if err != nil {
			return err
		}

		if flags.all {
			if len(args) > 0 {
				return fmt.Errorf("--all flag cannot be used with task IDs")
			}

			tasks, err := task.GetTasks()
			if err != nil {
				return fmt.Errorf("failed to retrieve tasks: %w", err)
			}

			removed := 0
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
					fmt.Printf("failed to remove task %d: %v\n", t.ID, err) // fixed message
				} else {
					removed++
				}
			}

			fmt.Printf("Removed %d tasks\n", removed)
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("task ID required")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("failed to parse task ID: %w", err)
		}

		if err := task.CheckTaskExists(id); err != nil {
			return fmt.Errorf("task does not exist: %w", err)
		}

		if err := task.RemoveTask(id); err != nil {
			return fmt.Errorf("failed to remove task: %w", err)
		}

		fmt.Println("Task removed")
		return nil
	},
}

func init() {
	removeCmd.Flags().BoolP("all", "a", false, "Remove all tasks (optionally with filters)")
	removeCmd.Flags().BoolP("priority", "p", false, "Remove only high-priority tasks (requires --all)")
	removeCmd.Flags().BoolP("normal", "n", false, "Remove only normal-priority tasks (requires --all)")
	removeCmd.Flags().BoolP("complete", "C", false, "Remove only complete tasks (optionally with filters)")
	removeCmd.Flags().BoolP("open", "o", false, "Remove only open tasks (optionally with filters)")

	rootCmd.AddCommand(removeCmd)
}
