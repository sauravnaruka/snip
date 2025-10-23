package search

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
)

func (c *Client) SearchMovie(query string) ([]Movie, error) {
	if query == "" {
		return nil, nil
	}

	tokens := preProcessQuery(query, c.stopWordMap)

	fmt.Printf("Tokens are: %s\n", tokens)

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

	return results, nil
}

func preProcessQuery(query string, stopWordMap map[string]struct{}) []string {
	queryLower := removePunctuation(query)
	tokens := createTokens(queryLower)
	tokens = processTokens(tokens, stopWordMap)
	return tokens
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

func processTokens(tokens []string, stopWordMap map[string]struct{}) []string {
	var result []string
	for _, token := range tokens {
		if !isStopWord(token, stopWordMap) {

			stemmed, err := snowball.Stem(token, "english", false)
			if err != nil {
				fmt.Printf("Error while converting word %s to stem.\n Error: %v", token, err)
				result = append(result, token)
				continue
			}

			result = append(result, stemmed)
		}
	}

	return result
}

func isStopWord(token string, stopWordMap map[string]struct{}) bool {
	_, ok := stopWordMap[token]
	return ok
}

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}

	return false
}
