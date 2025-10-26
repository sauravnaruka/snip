package search

import (
	"fmt"
	"sort"
)

type MovieBM25Score struct {
	ID    int
	Score float64
	Movie Movie
}

func (c *Client) SearchMovieBM25(query string, limit int) ([]MovieBM25Score, error) {
	if query == "" {
		return nil, nil
	}

	tokens := preProcessQuery(query, c.stopWordMap)
	fmt.Printf("User query tokens are: %s\n", tokens)

	docScoreMap := c.findMoviesByTokensBM25(tokens)
	if len(docScoreMap) == 0 {
		fmt.Println("No movies found.")
		return nil, nil
	}

	scores := make([]MovieBM25Score, 0, len(docScoreMap))
	for docID, score := range docScoreMap {
		movie, ok := c.invertedIndex.getMovieByID(docID)
		if !ok {
			fmt.Printf("Movie with ID=%d, not found", docID)
			continue
		}
		scores = append(scores, MovieBM25Score{docID, score, movie})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	// Get top N movies
	if limit > len(scores) {
		limit = len(scores)
	}

	results := scores[:limit]

	return results, nil
}

func (c *Client) findMoviesByTokensBM25(tokens []string) map[int]float64 {
	docScoreMap := make(map[int]float64)
	for _, token := range tokens {
		ids := c.invertedIndex.getDocumentIDs(token)

		for _, id := range ids {
			score, err := c.GetBM25(id, token)
			if err != nil {
				fmt.Printf("Error while calculating BM25 score for token='%s' in document with id=%v\n", token, id)
				continue
			}
			docScoreMap[id] += score
		}
	}

	return docScoreMap
}
