package hw03frequencyanalysis

import (
	"regexp"
	"sort"
)

type wordFrequency struct {
	word  string
	count int
}

var re = regexp.MustCompile(`\S+`)

func Top10(text string) []string {
	words := re.FindAllString(text, -1)

	wordCount := make(map[string]int)

	for _, word := range words {
		wordCount[word]++
	}

	frequencies := make([]wordFrequency, 0, len(wordCount))
	for word, count := range wordCount {
		frequencies = append(frequencies, wordFrequency{word, count})
	}

	sort.Slice(frequencies, func(i, j int) bool {
		if frequencies[i].count == frequencies[j].count {
			return frequencies[i].word < frequencies[j].word
		}
		return frequencies[i].count > frequencies[j].count
	})

	topWords := []string{}
	for i := 0; i < len(frequencies) && i < 10; i++ {
		topWords = append(topWords, frequencies[i].word)
	}

	return topWords
}
