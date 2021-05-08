package main

import (
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/tidwall/collate"
)

func main() {
	folder := "./words"
	prepareFolder(folder, "*.txt")

	for _, path := range getYamlFilepaths("../data/pl/") {
		meta := getMeta(path)
		filename := "słowa - " + meta.Label + ".txt"
		fmt.Println("Parsing...", filename)

		// set: extracted unique words normalized to lowercase
		set := make(map[string]void)
		extractWords(getRows(path).toString(), set)
		delete(set, "")

		// convert map[string]void to []string
		var words []string
		for word := range set {
			words = append(words, word)
		}

		sortArray(words, "POLISH_CI")

		var data []byte
		for _, word := range words {
			bytes := []byte(word + "\n")
			data = append(data, bytes...)
		}

		for err := ioutil.WriteFile(folder+"/"+filename, data, 0644); err != nil; {
			panic(err)
		}
	}
}

func sortArray(arr []string, lang string) {
	less := collate.IndexString(lang)
	sort.SliceStable(arr, func(i, j int) bool {
		return less(arr[i], arr[j])
	})
}