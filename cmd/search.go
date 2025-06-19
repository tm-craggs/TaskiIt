package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
	"github.com/tcraggs/TidyTask/util"
)

// create struct that defines the available flags for search command
type searchFlags struct {
	searchID       bool
	searchTitle    bool
	searchDue      bool
	filterComplete bool
	filterOpen     bool
	filterPriority bool
	filterNormal   bool
}

// helper function to parse flags with error handling
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

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [keyword] [flags]",
	Short: "Search for tasks using a keyword",
	Long: `The 'search' command is used to search for tasks using a keyword.

By default, the keyword is matched against all fields: ID, title, and due date. 
You can narrow the scope by explicitly specifying which fields to search using the --id, --title, and --due flags

"You can also narrow results using constraint flags, which show only tasks that meet the criteria you specify."`,

	Example: `  tidytask search essay
  > Search all fields for the word 'essay'

  tidytask search homework --title --priority
  > Search for tasks that contain 'homework' in the title, showing only priority results, 

  tidytask search 2024 --due --open --priority
  > Search due dates for the number '2024', show only tasks that are both open and high-priority`,

	RunE: func(cmd *cobra.Command, args []string) error {

		// check if keyword provided
		if len(args) == 0 {
			return fmt.Errorf("keyword required")
		}

		// get flags
		flags, err := getSearchFlags(cmd)
		if err != nil {
			return err
		}

		// if no fields selected, search all
		if !flags.searchID && !flags.searchTitle && !flags.searchDue {
			flags.searchID = true
			flags.searchTitle = true
			flags.searchDue = true
		}

		// get keyword
		keyword := args[0]

		// search specified fields for the keyword
		tasks, err := task.SearchTasks(keyword, flags.searchID, flags.searchTitle, flags.searchDue)
		if err != nil {
			return fmt.Errorf("failed searching tasks: %w", err)
		}

		// filter search results using flags
		filteredTasks := util.FilterTasks(tasks, flags.filterComplete, flags.filterPriority, flags.filterOpen, flags.filterNormal)

		// print tasks in table format
		util.PrintTasks(filteredTasks)

		// exit
		return nil
	},
}

// command initialisation
func init() {

	// define flags and add subcommand to root

	// target flags
	searchCmd.Flags().BoolP("id", "i", false, "Match keyword against task ID")
	searchCmd.Flags().BoolP("title", "t", false, "Match keyword against task title")
	searchCmd.Flags().BoolP("due", "d", false, "Match keyword against due date")

	// filter flags
	searchCmd.Flags().BoolP("complete", "c", false, "Show only complete tasks")
	searchCmd.Flags().BoolP("open", "o", false, "Show only open (incomplete) tasks")
	searchCmd.Flags().BoolP("priority", "p", false, "Show only high priority tasks")
	searchCmd.Flags().BoolP("normal", "n", false, "Show normal priority tasks")

	rootCmd.AddCommand(searchCmd)
}
