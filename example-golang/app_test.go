package main

import (
	"fmt"
	"os"
	"testing"
	"wordextractor/app"

	"github.com/stretchr/testify/assert"
)

func Test_ExtractUniqueWords(t *testing.T) {
	text := "ćma cześć ser. śmiech!żółw zebra-łuk len Ćma Żółw ser"
	expected := []string{"cześć", "ćma", "len", "łuk", "ser", "śmiech", "zebra", "żółw"}
	given, err := app.ExtractUniqueWords(text, "pl")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `ExtractUniqueWords error: %s`, err)
		return
	}
	assert.Equal(t, expected, given, "text should be tokenized into unique words")
}
