package search

import "fmt"

func (c *Client) GetTFIDF(docID int, token string) (float64, error) {
	tf, err := c.GetTermFrequency(docID, token)
	if err != nil {
		return 0, fmt.Errorf("error while getting term frequency with error: %w", err)
	}

	idf, err := c.GetInverseDocumentFrequency(token)
	if err != nil {
		return 0, fmt.Errorf("error while getting Inverse Document Frequency with error: %w", err)
	}

	tfIdf := float64(tf) * idf

	return tfIdf, nil
}
