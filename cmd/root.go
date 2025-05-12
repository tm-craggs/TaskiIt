/*
Copyright © 2025 NAME HERE <tom.craggs@protonmail.com>
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tidytask",
	Short: "A simple Linux CLI tool for managing your to-do-list. Built with Go",
	Long: `TidyTask is a CLI tool for managing your to-do-list. You can add, edit, list, search, complete and remove tasks
Tasks have an ID assigned by the program, and can have an optional message, deadline, and a priority status. These are
set and managed using flags.

Priority tasks are shown in blue and will always be at the top of your list. Tasks due today are shown in yellow,
overdue tasks in red, and completed tasks in green. 

Examples:
tidytask add --message "Finish Homework" --deadline 2025-08-15
-> create task "Finish Homework" with deadline August 15, 2025

tidytask remove 3 
-> remove task with ID 3 from system

tidytask search --message "Homework" 
-> search for all tasks containing "Homework" in their message)

tidytask edit 1 --deadline 2025-08-20
-> change deadline for task with ID 1 to August 20, 2025

tidytask complete 5 
-> mark task with ID 5 as complete

tidytask list 
-> list all tasks, closest deadline to furthest deadline

tidytask list --completed --reverse 
-> list all completed tasks in reverse order of completion

Use --help for more information.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.TidyTask.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
