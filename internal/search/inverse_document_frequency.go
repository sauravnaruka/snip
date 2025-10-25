package search

import (
	"fmt"
	"math"
)

func (c *Client) GetInverseDocumentFrequency(query string) (float64, error) {
	tokens := preProcessQuery(query, c.stopWordMap)
	if len(tokens) == 0 || len(tokens) > 1 {
		return 0, fmt.Errorf("term frequency can be searched for one token")
	}

	docIDs := c.invertedIndex.getDocumentIDs(tokens[0])
	termDocCount := len(docIDs)

	totalDocs := len(c.invertedIndex.DocMap)

	idf := math.Log(float64(totalDocs+1) / float64(termDocCount+1))
	return idf, nil
}
