package search

import "fmt"

func (c *Client) GetTermFrequency(docID int, query string) (int, error) {
	tokens := preProcessQuery(query, c.stopWordMap)
	if len(tokens) == 0 || len(tokens) > 1 {
		return 0, fmt.Errorf("term frequency can be searched for one token")
	}
	return c.invertedIndex.getTF(docID, tokens[0]), nil
}
