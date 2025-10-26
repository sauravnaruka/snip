package search

import "fmt"

func (c *Client) GetBM25(docID int, token string) (float64, error) {
	bm25TF, err := c.GetBM25TF(docID, token, c.BM25_K1, c.BM25_B)
	if err != nil {
		return 0.0, fmt.Errorf("error calculating BM25-TF. Error: %w", err)
	}

	bm25IDF, err := c.GetBM25IDF(token)
	if err != nil {
		return 0.0, fmt.Errorf("error calculating BM25-IDF. Error: %w", err)
	}

	bm25 := bm25TF * bm25IDF

	return bm25, nil
}
