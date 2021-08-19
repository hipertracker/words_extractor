package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_extractUniqueWords(t *testing.T) {
	text := "ćma cześć ser. śmiech!żółw zebra-łuk len Ćma Żółw ser"
	expected := []string{"ćma", "cześć", "ser", "śmiech", "żółw", "zebra", "łuk", "len"}
	given := extractUniqueWords([]byte(text))
	assert.Equal(t, expected, given, "text should be tokenized into unique words")
}

func Test_sortWords(t *testing.T) {
	words := []string{"ćma", "cześć", "ser", "śmiech", "żółw", "zebra", "łuk", "len"}
	expected := []string{"cześć", "ćma", "len", "łuk", "ser", "śmiech", "zebra", "żółw"}
	given := sortWords(words, "POLISH_CI")
	assert.Equal(t, expected, given, "words should be sorted out using Polish grammar rules")
}

// Total items: 123
// Total size: 503 MB`
// Total timing:  36.606038042s`
