package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_extractUniqueWords(t *testing.T) {
	text := "ćma cześć ser. śmiech!żółw zebra-łuk len Ćma Żółw ser"
	expected := []string{"ćma", "cześć", "ser", "śmiech", "żółw", "zebra", "łuk", "len"}
	given := extractUniqueWords([]byte(text))
	assert.Equal(t, expected, given, "text should be extracted into unique words")
}

func Test_sortWords(t *testing.T) {
	words := []string{"ćma", "cześć", "ser", "śmiech", "żółw", "zebra", "łuk", "len"}
	expected := []string{"cześć", "ćma", "len", "łuk", "ser", "śmiech", "zebra", "żółw"}
	given := sortWords(words, "POLISH_CI")
	assert.Equal(t, expected, given, "words should be sorted using Polish grammar rule")
}
