package search

import "fmt"

func (c *Client) GetBM25TF(docID int, query string) (float64, error) {
	tokens := preProcessQuery(query, c.stopWordMap)
	if len(tokens) == 0 || len(tokens) > 1 {
		return 0, fmt.Errorf("the method can calculate BM25-IDF for one token")
	}
	token := tokens[0]

	tf := c.invertedIndex.getTF(docID, token)

	bm25TF := (float64(tf) * (c.BM25_K1 + 1.0)) / (float64(tf) + c.BM25_K1)

	return bm25TF, nil
}
