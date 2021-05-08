package main

import (
	"os"
	"path/filepath"
	"strings"
)

type ListOfStrings []string
type void struct{}

var member void

func getYamlFilepaths(root string) []string {
	var result []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yml" {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return result
}

func removeCharacters(input string, characters string) string {
	filter := func(r rune) rune {
		if strings.IndexRune(characters, r) < 0 {
			return r
		}
		return -1
	}
	return strings.Map(filter, input)
}

func (r *resultsArray) extractWords(s string) {
	for _, word := range strings.Fields(s) {
		r.Results = append(r.Results, strings.ToLower(removeCharacters(word, ".:,;()!?'-_")))
	}
}

func prepareFolder(folder, pattern string) {
	os.Mkdir(folder, 0777)
	files, err := filepath.Glob(folder + "/" + pattern)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}
