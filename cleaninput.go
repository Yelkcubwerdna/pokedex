package main

import (
	"slices"
	"strings"
)

func cleanInput(text string) []string {
	var words []string

	text = strings.TrimSpace(text)

	words = strings.Split(text, " ")
	for i := 0; i < len(words); {

		words[i] = strings.ToLower(words[i])
		words[i] = strings.TrimSpace(words[i])

		if words[i] == "" {
			words = slices.Delete(words, i, i+1)
		} else {
			i++
		}
	}

	return words
}
