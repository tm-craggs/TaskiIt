/*
Copyright © 2025 NAME HERE <tom.craggs@protonmail.com>
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tidytask",
	Short: "A simple Linux CLI tool for managing your to-do-list. Built with Go",
	Long: `TidyTask CLI v1.0.0

A simple Linux CLI tool for managing your to-do-list. Built with Go

USAGE
$ tidytask COMMAND [FLAGS]

COMMANDS:
add		Create a new task
complete [ID]	Mark a task as complete by ID
remove [ID]	Remove a task by ID
edit [ID]	Edit a task by ID
list		List all tasks
search		Search tasks by keyword

Flags:
-h --help		Help for TidyTask

Use TidyTask [command] --help for more information about a command.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		task.InitDB()
	},
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
