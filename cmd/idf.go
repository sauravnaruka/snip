/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/SauravNaruka/snip/internal/search"
	"github.com/spf13/cobra"
)

func getIDFCmd(searchClient *search.Client) *cobra.Command {
	// searchCmd represents the search command
	return &cobra.Command{
		Use:   "idf <query>",
		Short: "Calculate the IDF for a given token",
		Long: `Calculate the Inverse Document Frequency (IDF) for a given token.
	
	Examples:
	  app idf "grizzly"
	  app idf actor`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			query := strings.Join(args, " ")

			idf, err := searchClient.GetInverseDocumentFrequency(query)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Printf("Inverse document frequency of '%s': %.2f\n", query, idf)
		},
	}
}
