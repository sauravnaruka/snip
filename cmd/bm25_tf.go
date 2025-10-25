package cmd

import (
	"fmt"
	"strconv"

	"github.com/SauravNaruka/snip/internal/search"
	"github.com/spf13/cobra"
)

func getBM25TFCmd(searchClient *search.Client) *cobra.Command {
	var k1 float64
	var b float64

	// searchCmd represents the search command
	cmd := &cobra.Command{
		Use:   "bm25tf <docID> <query>",
		Short: "Get BM25 TF score for a given term & document",
		Long: `Get BM25 TF score for a given term & document.
	
	Examples:
	  app bm25tf 1 anbuselvan
	  app bm25tf 1 maya --k1=1.5 --b=0.75`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Called with flags k1=%v, b=%v\n", k1, b)
			docID, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Invalid document ID: %s\n", args[0])
				return
			}

			token := args[1]

			bm25TF, err := searchClient.GetBM25TF(docID, token, k1, b)
			if err != nil {
				fmt.Printf("bm25tf faild with error: %v", err)
				return
			}

			fmt.Printf("BM25 TF score of '%s' in document '%v': '%.2f'\n", token, docID, bm25TF)
		},
	}

	cmd.Flags().Float64Var(&k1, "k1", searchClient.BM25_K1, "BM25 k1 parameter")
	cmd.Flags().Float64Var(&b, "b", searchClient.BM25_B, "BM25 b parameter is a tunable parameter that controls how much we care about document length.")

	return cmd
}
