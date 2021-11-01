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

	"github.com/bmatcuk/doublestar"
	"github.com/thoas/go-funk"
	"github.com/tidwall/collate"
)

type Pair struct {
	Path string
	Size int64
}

var srcPath = "../data/**/*.yml"
var outdir = "words"

func main() {
	paths, _ := doublestar.Glob(srcPath)

	clearResults()
	runWithChannels(paths)

	clearResults()
	runWithWaitGroups(paths)
}

func clearResults() {
	os.RemoveAll(outdir)
	os.Mkdir(outdir, 0777)
}

func runWithChannels(paths []string) {
	var ch = make(chan string)

	t := time.Now()
	defer timeTrack(t)

	for _, path := range paths {
		go func(yamlPath string) {
			ch <- parseFile(yamlPath, false)
		}(path)
	}
	for range paths {
		<-ch
	}
}

func runWithWaitGroups(paths []string) {
	var wg sync.WaitGroup
	t := time.Now()
	defer timeTrack(t)
	for _, path := range paths {
		wg.Add(1)
		go func(yamlPath string) {
			parseFile(yamlPath, false)
			wg.Done()
		}(path)
	}
	wg.Wait()
}

func parseFile(path string, sorting bool) string {
	// load YAML file
	meta := GetYAML(path)
	outfilepath := fmt.Sprintf("%s/extracted-words-for-%s.txt", outdir, meta.Code)

	// load text file
	filepath := strings.Replace(path, ".yml", ".txt", -1)
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	words := extractUniqueWords(content)

	// sort unique words
	if sorting {
		words = sortWords(words, "POLISH_CI")
	}

	text := strings.Join(words, "\n")
	for err := ioutil.WriteFile(outfilepath, []byte(text), 0644); err != nil; {
		panic(err)
	}
	return outfilepath
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
