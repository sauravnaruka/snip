package cmd

import (
	"fmt"
	"strings"

	"github.com/SauravNaruka/snip/internal/search"
	"github.com/spf13/cobra"
)

func getBM25IDFCmd(searchClient *search.Client) *cobra.Command {
	// searchCmd represents the search command
	return &cobra.Command{
		Use:   "bm25idf <query>",
		Short: "Get BM25 IDF score for a given term",
		Long: `Get BM25 IDF score for a given term.
	
	Examples:
	  app bm25idf "grizzly"
	  app bm25idf actor`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			query := strings.Join(args, " ")

			bm25IDF, err := searchClient.GetBM25IDF(query)
			if err != nil {
				fmt.Printf("GetBM25IDF faild with error: %v", err)
				return
			}

			fmt.Printf("BM25 IDF score of '%s': '%.2f'\n", query, bm25IDF)
		},
	}
}
