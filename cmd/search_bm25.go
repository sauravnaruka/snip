package cmd

import (
	"fmt"
	"strings"

	"github.com/SauravNaruka/snip/internal/search"
	"github.com/spf13/cobra"
)

func getSearchBM25Cmd(searchClient *search.Client) *cobra.Command {
	var limit int

	// searchCmd represents the search command
	cmd := &cobra.Command{
		Use:   "bm25search <query>",
		Short: "search for a movie or show using BM25",
		Long: `Search for a movie, show, or keyword in the Snip catalog using BM25.
	
	Examples:
	  app bm25search "Inception"
	  app bm25search Stranger`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			searchStr := strings.Join(args, " ")
			fmt.Printf("Searching for: %s\n", searchStr)

			results, err := searchClient.SearchMovieBM25(searchStr, limit)
			if err != nil {
				fmt.Printf("Search faild with error: %v", err)
				return
			}

			for i, r := range results {
				fmt.Printf("%d. (%d) %s - Score: %.2f\n", i+1, r.ID, r.Movie.Title, r.Score)
			}
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 5, "limit parameter limits the number of results")

	return cmd
}
