package main

import (
	"os"
	"path/filepath"
	"regexp"
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

func extractWords(s string, set map[string]void) {
	re := regexp.MustCompile("[^\\p{L}]+")
	for _, word := range re.Split(s, -1) {
		set[strings.ToLower(word)] = member
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
