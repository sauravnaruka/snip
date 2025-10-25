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

func getTFIDFCmd(searchClient *search.Client) *cobra.Command {
	// searchCmd represents the search command
	return &cobra.Command{
		Use:   "tfidf <docID> <query>",
		Short: "Calculate the TF-IDF for a given token in a document",
		Long: `Calculate the Term Frequency and Inverse Document Frequency (TF-IDF) for a given token.
	Formula:
		TF-IDF = TF * IDF

	Examples:
	  app idf "grizzly"
	  app idf actor`,

		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			docID, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Invalid document ID: %s\n", args[0])
				return
			}

			token := args[1]

			tfIdf, err := searchClient.GetTFIDF(docID, token)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Printf("TF-IDF score of '%s' in document '%v' %.2f\n", token, docID, tfIdf)
		},
	}
}
