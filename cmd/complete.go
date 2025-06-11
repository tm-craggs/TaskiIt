package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"strconv"
)

type completeFlags struct {
	all      bool
	priority bool
	normal   bool
}

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

// completeCmd represents the `complete` command for marking tasks as done
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Complete tasks",
	Long:  `Mark tasks as complete, by ID or using filters`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := task.BackupDB(); err != nil {
			return fmt.Errorf("failed to back up database: %w", err)
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := getCompleteFlags(cmd)
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

			for _, t := range tasks {
				if flags.priority && !t.Priority {
					continue
				}
				if flags.normal && t.Priority {
					continue
				}
				if err := task.CompleteTask(t.ID); err != nil {
					fmt.Printf("failed to complete task %d: %v\n", t.ID, err)
				}
			}

			fmt.Println("Filtered tasks marked complete")
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
	completeCmd.Flags().Bool("all", false, "Complete all tasks (optionally with filters)")
	completeCmd.Flags().Bool("priority", false, "Complete only high-priority tasks (requires --all)")
	completeCmd.Flags().Bool("normal", false, "Complete only normal-priority tasks (requires --all)")
	rootCmd.AddCommand(completeCmd)
}
