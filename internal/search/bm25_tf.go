package search

import "fmt"

func (c *Client) GetBM25TF(docID int, query string, k1 float64, b float64) (float64, error) {
	tokens := preProcessQuery(query, c.stopWordMap)
	if len(tokens) == 0 || len(tokens) > 1 {
		return 0, fmt.Errorf("the method can calculate BM25-IDF for one token")
	}
	token := tokens[0]

	tf := c.invertedIndex.getTF(docID, token)

	lengthNorm := c.getBM25LengthNormalizer(docID)

	bm25TF := (float64(tf) * (c.BM25_K1 + 1.0)) / (float64(tf) + c.BM25_K1*lengthNorm)

	return bm25TF, nil
}

func (c *Client) getBM25LengthNormalizer(docID int) float64 {
	docLength, ok := c.invertedIndex.DocLengthMap[docID]
	if !ok {
		fmt.Printf("Warning: document with ID=%v has no entry. Taking document length as 0", docID)
		docLength = 1.0
	}

	avgDocLength := c.invertedIndex.getAvgDocLength()
	if avgDocLength == 0 {
		// No documents in index, return neutral normalization factor
		return 1.0
	}

	return 1 - c.BM25_B + c.BM25_B*(float64(docLength)/avgDocLength)
}
