package cmd

import (
	"fmt"
	"strconv"

	"github.com/SauravNaruka/snip/internal/search"
	"github.com/spf13/cobra"
)

func getBM25TFCmd(searchClient *search.Client) *cobra.Command {
	// searchCmd represents the search command
	return &cobra.Command{
		Use:   "bm25tf <docID> <query>",
		Short: "Get BM25 TF score for a given term & document",
		Long: `Get BM25 TF score for a given term & document.
	
	Examples:
	  app bm25tf 1 anbuselvan
	  app bm25tf 1 maya`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			docID, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Invalid document ID: %s\n", args[0])
				return
			}

			token := args[1]

			bm25TF, err := searchClient.GetBM25TF(docID, token)
			if err != nil {
				fmt.Printf("bm25tf faild with error: %v", err)
				return
			}

			fmt.Printf("BM25 TF score of '%s' in document '%v': '%.2f'\n", token, docID, bm25TF)
		},
	}
}
