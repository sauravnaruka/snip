package search

import (
	"fmt"
	"strings"
)

func (idx *Client) SearchMovie(query string) ([]Movie, error) {
	if query == "" {
		return nil, nil
	}

	tokens := preProcessQuery(query, idx.stopWordMap)

	fmt.Printf("Tokens are: %s\n", tokens)

	var results []Movie
	for _, movie := range idx.movies {
		titleClean := removePunctuation(movie.Title)

		if containsAny(titleClean, tokens) {
			results = append(results, movie)
		}
	}

	if len(results) == 0 {
		fmt.Println("No movies found.")
	}

	return results, nil
}

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}

	return false
}
