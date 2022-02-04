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
	SrcPath string
	Dstpath string
}

var srcPath = "../data/??/**/*.yml"
var outdir = "words"

var wg sync.WaitGroup

// func mainOld() {
// 	paths, _ := doublestar.Glob(srcPath)

// 	clearResults()
// 	runWithChannels(paths)

// 	clearResults()
// 	runWithWaitGroups(paths)
// }

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		t := time.Now()
		defer timeTrack(t)

		paths, _ := doublestar.Glob(srcPath)

		ch1 := make(chan Pair, len(paths))
		ch2 := make(chan string, len(paths))

		clearResults()

		for _, yamlPath := range paths {
			go loadYaml(ch1, yamlPath)
		}

		for range paths {
			pair := <-ch1
			go loadText(ch2, pair.SrcPath, pair.Dstpath, true)
		}
		for range paths {
			fmt.Printf("Saved %s\n", <-ch2)
		}
	}()
	wg.Wait()
}

func loadYaml(ch chan Pair, path string) {
	meta := GetYAML(path)
	srcPath := strings.Replace(path, ".yml", ".txt", -1)
	dstPath := fmt.Sprintf("%s/extracted-words-for-%s.txt", outdir, meta.Code)
	ch <- Pair{srcPath, dstPath}
}

func loadText(ch2 chan string, srcPath string, dstPath string, sorting bool) {
	content, err := ioutil.ReadFile(srcPath)
	if err != nil {
		panic(err)
	}
	words := extractUniqueWords(content)
	if sorting {
		words = sortWords(words, "POLISH_CI")
	}
	text := strings.Join(words, "\n")
	for err := ioutil.WriteFile(dstPath, []byte(text), 0644); err != nil; {
		panic(err)
	}
	ch2 <- dstPath
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
