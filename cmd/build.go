package cmd

import (
	"fmt"

	"github.com/SauravNaruka/snip/internal/search"
	"github.com/spf13/cobra"
)

func getBuildCmd(searchClient *search.Client) *cobra.Command {
	// searchCmd represents the search command
	return &cobra.Command{
		Use:   "build",
		Short: "Build inverted index for the movie database",
		Long: `Build inverted index for the movie database.
	
	Examples:
	  app build`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Creating inverted index for data")

			err := searchClient.BuildInvertedIndex()
			if err != nil {
				fmt.Printf("Failed to build inverted index with error: %v", err)
				return
			}

		},
	}
}
