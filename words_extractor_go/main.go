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

type Pair struct {
	Path string
	Size int64
}

func main() {
	// var wg sync.WaitGroup
	queue := make(chan Pair)

	t := time.Now()
	defer timeTrack(t)

	outdir := "words"
	os.RemoveAll(outdir)
	os.Mkdir(outdir, 0777)

	fmt.Println("Parsing...")

	paths, _ := doublestar.Glob("../data/**/*.yml")

	total_size := int64(0)
	items_count := len(paths)
	for i, path := range paths {
		go processFile(queue, outdir, path, false)
		res := <-queue
		total_size += res.Size
		fmt.Printf("[%d/%d] %s\n", i+1, items_count, res.Path)
	}
	fmt.Printf("Total items: %d\n", items_count)
	fmt.Printf("Total size: %d MB\n", total_size/(1024*1024))
}

func processFile(queue chan Pair, outdir string, path string, sorting bool) {
	meta := GetYAML(path)
	// load text file
	filepath := strings.Replace(path, ".yml", ".txt", -1)
	info, err := os.Stat(filepath)
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	// extract and sort unique words
	words := extractUniqueWords(content)
	if sorting {
		words = sortWords(words, "POLISH_CI")
	}
	text := strings.Join(words, "\n")
	outpath := fmt.Sprintf("%s/%s-%s.txt", outdir, meta.Lang, meta.Code)
	for err := ioutil.WriteFile(outpath, []byte(text), 0644); err != nil; {
		panic(err)
	}
	queue <- Pair{path, info.Size()}
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
