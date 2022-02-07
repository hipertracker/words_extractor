package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/bmatcuk/doublestar"
)

const dirPerms = 0755

type empty struct{}

// Run extracts unique words from the list of files and saves them to the outDir.
// No error handling, no context cancellation is implemented to match implementations
// in other languages.
func Run(srcDir, outDir string, numWorkers int, sortResults bool) error {
	files, err := doublestar.Glob(srcDir)
	if err != nil {
		return fmt.Errorf(`app: getting list of files "%s": %w`, srcDir, err)
	}

	if err = clearOutput(outDir); err != nil {
		return err
	}

	// This is a very basic semaphore implementation. Counting unique words from
	// a stream of data is IO, memory and CPU expensive. Semaphore lets to run
	// up to the numWorkers or workers concurrently and, by default, this number
	// matches the number of CPUs.
	sem := make(chan empty, numWorkers)

	var wg sync.WaitGroup
	var spec *MetaConfig

	for _, file := range files {
		sem <- empty{}

		if spec, err = ReadSpec(file); err != nil {
			return err
		}

		src := file[:len(file)-3] + "txt"
		dst := filepath.Join(outDir, spec.Lang+"-"+spec.Code+".txt")

		wg.Add(1)
		go extract(src, dst, langMap[spec.Lang], sortResults, sem, &wg)
	}

	wg.Wait()
	close(sem)

	return nil
}

func clearOutput(outDir string) error {
	if err := os.RemoveAll(outDir); err != nil {
		return fmt.Errorf(`app: cleaning previous results in "%s": %w`, outDir, err)
	}
	if err := os.MkdirAll(outDir, dirPerms); err != nil {
		return fmt.Errorf(`app: creating output directory "%s": %w`, outDir, err)
	}

	return nil
}

var langMap = map[string]string{
	"en": "ENGLISH_CI", // The first language is used as fallback.
	"la": "ENGLISH_CI", // Latin
	"eo": "ENGLISH_CI", // Esperanto
	"ar": "ARABIC_CI",
	"cz": "CZECH_CI",
	"da": "DANISH_CI", // ?
	"de": "GERMAN_CI",
	"el": "GREEK_CI",
	"es": "SPANISH_CI",
	"fi": "FINNISH_CI",
	"fr": "FRENCH_CI",
	"he": "HEBREW_CI",
	"hr": "CROATIAN_CI",
	"hu": "HUNGARIAN_CI",
	"it": "ITALIAN_CI",
	"lt": "LITHUANIAN_CI",
	"nl": "DUTCH_CI",
	"pl": "POLISH_CI",
	"pt": "PORTUGUESE_CI",
	"ru": "RUSSIAN_CI",
	"sk": "SLOVAK_CI",
	"sv": "SWEDISH_CI",
	"uk": "UKRAINIAN_CI",
}
