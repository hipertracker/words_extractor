package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v2"
	"github.com/thoas/go-funk"
	"github.com/tidwall/collate"
)

func main() {
	t := time.Now()
	defer timeTrack(t)

	outdir := "words"
	os.RemoveAll(outdir)
	os.Mkdir(outdir, 0777)

	paths, _ := doublestar.Glob("../data/pl/**/*.yml")
	for _, path := range paths {
		fmt.Println(filepath.Base(path))
		meta := GetYAML(path)
		outfilepath := fmt.Sprintf("%s/extracted-words-for-%s.txt", outdir, meta.Code)
		
		// load file content
		filepath := strings.Replace(path, ".yml", ".txt", -1)
		content, _ := ioutil.ReadFile(filepath)

		// extract and sort unique words
		words := extractUniqueWords(content)
		words = sortWords(words, "POLISH_CI")

		text := strings.Join(words, "\n")
		for err := ioutil.WriteFile(outfilepath, []byte(text), 0644); err != nil; {
			panic(err)
		}
	}

}

func timeTrack(start time.Time) {
	fmt.Println("Total timing: ", time.Since(start))
}

func extractUniqueWords(content []byte) []string {
	text := strings.ToLower(string(content))
	re := regexp.MustCompile(`[^\p{L}]+`)
	tokens := re.Split(text, -1)
	return funk.UniqString(tokens)
}

func sortWords(words []string, lang string) []string {
	less := collate.IndexString(lang)
	sort.SliceStable(words, func(i, j int) bool {
		return less(words[i], words[j])
	})
	return words
}
