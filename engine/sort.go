package engine

import (
	"sort"
)

func SortResult(matchedIDs []uint64, docs []Documents, query string) []uint64 {
	sortMap := make(map[uint64]float64)
	for _, id := range matchedIDs {
		doc := docs[id]
		var words []string
		for _, token := range analyze(doc.Abstract) {
			words = append(words, token)
			if query == token {
			}
		}
		sortMap[id] = TFIDF(words, query, uint64(len(docs)), uint64(len(matchedIDs)))
	}

	sort.Slice(matchedIDs, func(i, j int) bool {
		return sortMap[matchedIDs[i]] > sortMap[matchedIDs[j]]
	})

	return matchedIDs
}
