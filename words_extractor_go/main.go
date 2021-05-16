package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bmatcuk/doublestar/v2"
	"github.com/thoas/go-funk"
	"github.com/tidwall/collate"
)

func main() {
	var wg sync.WaitGroup

	t := time.Now()
	defer timeTrack(t)

	outdir := "words"
	os.RemoveAll(outdir)
	os.Mkdir(outdir, 0777)

	fmt.Println("Parsing...")

	paths, _ := doublestar.Glob("../data/pl/**/*.yml")
	for i, path := range paths {
		wg.Add(1)
		go worker(i, &wg, path, outdir, true)
	}
	wg.Wait()
}

func worker(id int, wg *sync.WaitGroup, path, outdir string, verbose bool) {
	defer wg.Done()
	// load YAML file
	meta := GetYAML(path)
	outfilepath := fmt.Sprintf("%s/extracted-words-for-%s.txt", outdir, meta.Code)

	// load text file
	filepath := strings.Replace(path, ".yml", ".txt", -1)
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	// extract and sort unique words
	words := extractUniqueWords(content)
	words = sortWords(words, "POLISH_CI")

	text := strings.Join(words, "\n")
	for err := ioutil.WriteFile(outfilepath, []byte(text), 0644); err != nil; {
		panic(err)
	}
	if verbose {
		fmt.Println("Saved ", path)
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
