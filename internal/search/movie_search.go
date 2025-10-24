package search

import (
	"fmt"
)

func (c *Client) SearchMovie(query string) ([]Movie, error) {
	if query == "" {
		return nil, nil
	}

	tokens := preProcessQuery(query, c.stopWordMap)
	fmt.Printf("User query tokens are: %s\n", tokens)

	c.invertedIndex.load()

	results := c.findMoviesByTokens(tokens)

	if len(results) == 0 {
		fmt.Println("No movies found.")
	}

	return results, nil
}

func (c *Client) findMoviesByTokens(tokens []string) []Movie {
	docIDs := make(map[int]struct{})

outerLoop:
	for _, token := range tokens {
		ids := c.invertedIndex.getDocumentIDs(token)

		for _, id := range ids {
			docIDs[id] = struct{}{}

			if len(docIDs) == 5 {
				break outerLoop
			}
		}
	}

	results := make([]Movie, 0, len(docIDs))
	for id := range docIDs {
		m, ok := c.invertedIndex.getMovieByID(id)
		if !ok {
			fmt.Printf("Warning: movie with id=%v, not found\n", id)
			continue
		}

		results = append(results, m)
	}

	return results

}
