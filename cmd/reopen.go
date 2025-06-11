package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"strconv"
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
	Long:  `Remove a tasks complete status, by ID or using filters`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := task.BackupDB(); err != nil {
			return fmt.Errorf("failed to back up database: %w", err)
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := getReopenFlags(cmd)
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

			reopened := 0
			for _, t := range tasks {
				if flags.priority && !t.Priority {
					continue
				}
				if flags.normal && t.Priority {
					continue
				}

				if !t.Complete {
					continue
				}

				if err := task.ReopenTask(t.ID); err != nil {
					fmt.Printf("failed to complete task %d: %v\n", t.ID, err)
				} else {
					reopened++
				}
			}

			fmt.Printf("Reopened %d tasks\n", reopened)
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

		if err := task.CompleteTask(id); err != nil {
			return fmt.Errorf("failed to complete task: %w", err)
		}

		fmt.Println("Task completed successfully")
		return nil
	},
}

func init() {
	reopenCmd.Flags().BoolP("all", "a", false, "Complete all tasks (optionally with filters)")
	reopenCmd.Flags().BoolP("priority", "p", false, "Complete only high-priority tasks (requires --all)")
	reopenCmd.Flags().BoolP("normal", "n", false, "Complete only normal-priority tasks (requires --all)")

	rootCmd.AddCommand(reopenCmd)
}
