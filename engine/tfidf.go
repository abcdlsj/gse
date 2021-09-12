package engine

import (
	"math"
)

func TermFrequency(doc []string, query string) float64 {
	var maxFreq int
	words := make(map[string]int)

	for _, word := range doc {
		words[word]++
		if words[word] > maxFreq {
			maxFreq = words[word]
		}
	}
	return 0.5 * (1 + float64(words[query])/float64(maxFreq))
}

func InverseDocumentFrequency(docLen, curLen uint64) float64 {
	return math.Log(float64(docLen)) - math.Log(float64(curLen)+1)
}

func TFIDF(docWords []string, queryStr string, docLen, curLen uint64) float64 {
	return TermFrequency(docWords, queryStr) * InverseDocumentFrequency(docLen, curLen)
}
