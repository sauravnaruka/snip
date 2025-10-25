package search

import (
	"fmt"
	"math"
)

func (c *Client) GetBM25IDF(query string) (float64, error) {
	tokens := preProcessQuery(query, c.stopWordMap)
	if len(tokens) == 0 || len(tokens) > 1 {
		return 0, fmt.Errorf("the method can calculate BM25-IDF for one token")
	}
	token := tokens[0]

	totalDocs := len(c.invertedIndex.DocMap)

	docIDs := c.invertedIndex.getDocumentIDs(token)
	termDocCount := len(docIDs)

	docsWithoutToken := totalDocs - termDocCount

	return math.Log((float64(docsWithoutToken)+.5)/(float64(termDocCount)+.5) + 1), nil
}
