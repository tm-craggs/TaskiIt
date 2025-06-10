package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
)

type searchFlags struct {
	searchID       bool
	searchTitle    bool
	searchDue      bool
	filterComplete bool
	filterOpen     bool
	filterPriority bool
	filterNormal   bool
}

func getSearchFlags(cmd *cobra.Command) (*searchFlags, error) {
	flags := &searchFlags{}
	var err error

	if flags.searchID, err = cmd.Flags().GetBool("id"); err != nil {
		return nil, fmt.Errorf("failed to parse --id flag: %w", err)
	}
	if flags.searchTitle, err = cmd.Flags().GetBool("title"); err != nil {
		return nil, fmt.Errorf("failed to parse --title flag: %w", err)
	}
	if flags.searchDue, err = cmd.Flags().GetBool("due"); err != nil {
		return nil, fmt.Errorf("failed to parse --due flag: %w", err)
	}
	if flags.filterComplete, err = cmd.Flags().GetBool("complete"); err != nil {
		return nil, fmt.Errorf("failed to parse --complete flag: %w", err)
	}
	if flags.filterOpen, err = cmd.Flags().GetBool("open"); err != nil {
		return nil, fmt.Errorf("failed to parse --open flag: %w", err)
	}
	if flags.filterPriority, err = cmd.Flags().GetBool("priority"); err != nil {
		return nil, fmt.Errorf("failed to parse --priority flag: %w", err)
	}
	if flags.filterNormal, err = cmd.Flags().GetBool("normal"); err != nil {
		return nil, fmt.Errorf("failed to parse --normal flag: %w", err)
	}

	return flags, nil
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Add all new task to your to-do list",
	Long:  `Long description goes here`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no search word specified")
		}

		flags, err := getSearchFlags(cmd)
		if err != nil {
			return err
		}

		keyword := args[0]

		tasks, err := task.SearchTasks(keyword, flags.searchID, flags.searchTitle, flags.searchDue)
		if err != nil {
			return fmt.Errorf("failed searching tasks: %w", err)
		}

		filteredTasks := util.FilterTasks(tasks, flags.filterComplete, flags.filterPriority, flags.filterOpen, flags.filterNormal)

		util.PrintTasks(filteredTasks)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// target flags
	searchCmd.Flags().BoolP("id", "i", false, "Search task by ID")
	searchCmd.Flags().BoolP("title", "t", false, "Search task by title")
	searchCmd.Flags().BoolP("due", "d", false, "Search task by due")

	// filter flags
	searchCmd.Flags().BoolP("complete", "c", false, "Search only complete tasks")
	searchCmd.Flags().BoolP("open", "o", false, "Search only open tasks")
	searchCmd.Flags().BoolP("priority", "p", false, "Search task by priority")
	searchCmd.Flags().BoolP("normal", "n", false, "Search normal priority tasks")
}
