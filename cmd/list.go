package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/tidytask/task"
	"github.com/tcraggs/tidytask/util"
)

// create struct that defines the available flags for list command
type listFlags struct {
	priority bool
	complete bool
	open     bool
	normal   bool
}

// helper function to parse flags with error handling
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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display tasks in your to-do list",
	Long: `The 'list' command displays all tasks in your to-do list. 

Optionally, you can use flags to to narrow the results and only show tasks that meet certain criteria.`,

	Example: `  tidytask list
  > Show all tasks
  
  tidytask list --priority
  > Show only high priority tasks

  tidytask list --complete --priority
  > Show only completed, high priority tasks`,

	RunE: func(cmd *cobra.Command, args []string) error {

		// check args
		if len(args) > 0 {
			return fmt.Errorf("unexpected arguments: %v; use --help for usage information", args)
		}

		// get flags
		flags, err := getListFlags(cmd)
		if err != nil {
			return err
		}

		// check for flag conflicts
		if flags.priority && flags.normal {
			return fmt.Errorf("conflicting flags: cannot use --priority and --normal together")
		}

		if flags.complete && flags.open {
			return fmt.Errorf("conflicting flags: cannot use --complete and --open together")
		}

		// get tasks
		tasks, err := task.GetTasks()
		if err != nil {
			return fmt.Errorf("failed to get tasks: %w", err)
		}

		// filter tasks using flags
		filteredTasks := util.FilterTasks(tasks, flags.complete, flags.priority, flags.open, flags.normal)

		// print tasks in table format
		util.PrintTasks(filteredTasks)

		// exit
		return nil
	},
}

// command initialisation
func init() {

	// define flags and add subcommand to root

	listCmd.Flags().BoolP("priority", "p", false, "Show only high priority tasks")
	listCmd.Flags().BoolP("complete", "c", false, "Show only complete tasks ")
	listCmd.Flags().BoolP("open", "o", false, "Show only open (incomplete) tasks")
	listCmd.Flags().BoolP("normal", "n", false, "Show only normal priority tasks")

	rootCmd.AddCommand(listCmd)
}
