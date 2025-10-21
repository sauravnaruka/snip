package search

import (
	"fmt"
	"strings"
	"unicode"
)

func (c *Client) SearchMovie(query string) ([]Movie, error) {
	if query == "" {
		return nil, nil
	}

	queryLower := removePunctuation(query)
	tokens := createTokens(queryLower)
	tokens = removeStopWords(tokens, c.stopWordMap)

	fmt.Printf("Tokens are: %s", tokens)

	var results []Movie
	for _, movie := range c.movies {
		titleClean := removePunctuation(movie.Title)

		if containsAny(titleClean, tokens) {
			results = append(results, movie)
		}
	}

	if len(results) == 0 {
		fmt.Println("No movies found.")
	}

	for i, movie := range results {
		fmt.Printf("%d. %s\n", i+1, movie.Title)
	}

	return results, nil
}

func removePunctuation(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

func createTokens(s string) []string {
	tokens := strings.Fields(s)

	return tokens
}

func removeStopWords(tokens []string, dict map[string]string) []string {
	var result []string
	for _, t := range tokens {
		if _, ok := dict[t]; !ok {
			result = append(result, t)
		}
	}

	return result
}

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}

	return false
}
