package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
)

type listFlags struct {
	priority bool
	complete bool
	open     bool
	normal   bool
}

func getListFlags(cmd *cobra.Command) (listFlags, error) {
	var flags listFlags
	var err error

	if flags.priority, err = cmd.Flags().GetBool("priority"); err != nil {
		return flags, fmt.Errorf("failed to parse --priority flag: %w", err)
	}
	if flags.complete, err = cmd.Flags().GetBool("complete"); err != nil {
		return flags, fmt.Errorf("failed to parse --complete flag: %w", err)
	}
	if flags.open, err = cmd.Flags().GetBool("open"); err != nil {
		return flags, fmt.Errorf("failed to parse --open flag: %w", err)
	}
	if flags.normal, err = cmd.Flags().GetBool("normal"); err != nil {
		return flags, fmt.Errorf("failed to parse --normal flag: %w", err)
	}

	return flags, nil
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `List all tasks`,

	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := getListFlags(cmd)
		if err != nil {
			return err
		}

		tasks, err := task.GetTasks()
		if err != nil {
			return fmt.Errorf("failed to get tasks: %w", err)
		}

		util.PrintTasks(util.FilterTasks(tasks, flags.complete, flags.priority, flags.open, flags.normal))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("priority", "p", false, "Filter by priority")
	listCmd.Flags().BoolP("complete", "c", false, "Filter by complete")
	listCmd.Flags().BoolP("open", "o", false, "Filter by open task")
	listCmd.Flags().BoolP("normal", "n", false, "Filter by normal priority")
}
