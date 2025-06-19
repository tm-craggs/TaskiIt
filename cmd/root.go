/*
Copyright © 2025 NAME HERE <tom.craggs@protonmail.com>
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tcraggs/TidyTask/task"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tidytask",
	Short:   "A simple Linux CLI tool for managing your to-do-list. Built with Go",
	Version: "1.0.0",
	Long: `
 _____ _     _     _____         _    
|_   _(_)   | |   |_   _|       | |   
  | |  _  __| |_   _| | __ _ ___| | __
  | | | |/ _` + "`" + ` | | | | |/ _` + "`" + ` / __| |/ /
  | | | | (_| | |_| | | (_| \__ \   < 
  \_/ |_|\__,_|\__, \_/\__,_|___/_|\_\
                __/ |                 
               |___/                  

A simple CLI tool for managing your to-do-list. Built with Go`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := task.InitDB(); err != nil {
			return fmt.Errorf("DB creation error: %w", err)
		}
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if err := task.CloseDB(); err != nil {
			return fmt.Errorf("DB closing error: %v\n", err)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SilenceUsage = true                     // Show usage when error occurs
	rootCmd.SilenceErrors = true                    // Print errors
	rootCmd.TraverseChildren = false                // Don't defer flags to subcommands
	rootCmd.FParseErrWhitelist.UnknownFlags = false // Be strict about flags

	if err := rootCmd.Execute(); err != nil {
		// This prints the error to stderr
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
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
}
