package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	dict, keys := createDictionary(words)
	sort.Slice(keys, freqSorter(dict, keys))

	length := min(10, len(keys))
	return keys[:length]
}

func freqSorter(dict map[string]int, keys []string) func(i int, j int) bool {
	return func(i, j int) bool {
		if dict[keys[i]] != dict[keys[j]] {
			return dict[keys[i]] > dict[keys[j]]
		}
		return keys[i] < keys[j]
	}
}

func createDictionary(words []string) (map[string]int, []string) {
	dict := make(map[string]int)
	keys := make([]string, 0)
	for _, word := range words {
		trimmedWord := strings.Trim(strings.ToLower(word), ".!?,:`'\"")
		if trimmedWord == "" || trimmedWord == "-" {
			continue
		}
		curCount, ok := dict[trimmedWord]
		if ok {
			dict[trimmedWord] = curCount + 1
		} else {
			keys = append(keys, trimmedWord)
			dict[trimmedWord] = 1
		}
	}
	return dict, keys
}
