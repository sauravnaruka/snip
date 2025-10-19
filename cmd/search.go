/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/SauravNaruka/snip/internal/search"
	"github.com/spf13/cobra"
)

func newSearchCmd(searchClient *search.Client) *cobra.Command {
	// searchCmd represents the search command
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search for a movie or show",
		Long: `Search for a movie, show, or keyword in the Snip catalog.
	
	Examples:
	  app search "Inception"
	  app search Stranger`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			searchStr := args[0]
			fmt.Printf("Searching for: %s\n", searchStr)
			searchClient.SearchMovie(searchStr)
		},
	}
}
