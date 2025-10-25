/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/SauravNaruka/snip/internal/search"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

func GetRootCommand(c *search.Client) *cobra.Command {

	if rootCmd != nil {
		return rootCmd
	}

	rootCmd := &cobra.Command{
		Use:   "snip",
		Short: "Snip is a search tool for internal database",
		Long:  `Snip uses RAG (Retrieval Augmented Generation).`,
		Run: func(cmd *cobra.Command, args []string) {
			// This runs when no subcommand is provided
			fmt.Println("Available Commands:")
			cmd.Help()
		},
	}

	rootCmd.AddCommand(getSearchCmd(c))
	rootCmd.AddCommand(getBuildCmd(c))
	rootCmd.AddCommand(getTFCmd(c))
	rootCmd.AddCommand(getIDFCmd(c))

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	if rootCmd == nil {
		fmt.Println("Executing command without initializing cmd")
	}

	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
