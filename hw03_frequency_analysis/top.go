package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(s string) []string {
	if len(s) == 0 {
		return []string{}
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{}
	}
	frequency := make(map[string]int)
	for _, w := range words {
		frequency[w]++
	}
	uniqueWords := make([]string, 0)
	for w := range frequency {
		uniqueWords = append(uniqueWords, w)
	}
	sort.Slice(
		uniqueWords,
		func(i, j int) bool {
			diff := frequency[uniqueWords[i]] - frequency[uniqueWords[j]]
			if diff == 0 {
				return strings.Compare(uniqueWords[i], uniqueWords[j]) < 0
			}
			return diff > 0
		})
	topLimit := 10
	if len(uniqueWords) <= topLimit {
		return uniqueWords
	}
	return uniqueWords[:topLimit]
}
