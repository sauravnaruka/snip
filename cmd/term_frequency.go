/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/SauravNaruka/snip/internal/search"
	"github.com/spf13/cobra"
)

func getTFCmd(searchClient *search.Client) *cobra.Command {
	// searchCmd represents the search command
	return &cobra.Command{
		Use:   "tf <docID> <query>",
		Short: "Find term frequency for a token in a given document",
		Long: `Find term frequency for a token in a given document.
	
	Examples:
	  app tf 424 "trapper"
	  app tf 424 bear`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			docID, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Invalid document ID: %s\n", args[0])
				return
			}

			token := args[1]
			freq, err := searchClient.GetTermFrequency(docID, token)
			if err != nil {
				fmt.Printf("GetTermFrequency faild with error: %v", err)
				return
			}

			fmt.Printf("freq = %v\n", freq)
		},
	}
}
