package app

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/tidwall/collate"
)

const filePerm = 0644

// splitWordsUnicode splits data into words, using Unicode Letter character class.
// It works similar to the regular expression "[^\p{L}]+". This is what was used
// in the original code. Unicode function has slight overhead, but handles UTF-8
// correctly.
//
// Rust and Python versions split text according to "[\W\d]+" - anything that is
// not a word or a digit. WTF?
func splitWordsUnicode(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var start int
	var r rune
	for width := 0; start < len(data); start += width {
		if r, width = utf8.DecodeRune(data[start:]); unicode.IsLetter(r) {
			break
		}
	}

	for width, i := 0, start; i < len(data); i += width {
		if r, width = utf8.DecodeRune(data[i:]); !unicode.IsLetter(r) {
			return i + width, data[start:i], nil
		}
	}

	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	return start, nil, nil
}

// splitWords splits data into words similar to the "[\W\d]+" regular expression.
func splitWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var start int
	var r rune
	for width := 0; start < len(data); start += width {
		if r, width = utf8.DecodeRune(data[start:]); isLatin(r) {
			break
		}
	}

	for width, i := 0, start; i < len(data); i += width {
		if r, width = utf8.DecodeRune(data[i:]); !isLatin(r) {
			return i + width, data[start:i], nil
		}
	}

	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	return start, nil, nil
}

func isLatin(r rune) bool {
	if r >= 0x80 || r == 0x00 {
		return false
	}

	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}

func extract(src, dst, lang string, sortResults bool, sem <-chan empty, wg *sync.WaitGroup) {
	defer func() {
		<-sem
		wg.Done()
	}()

	fd, err := os.Open(src)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `extract: opening source file "%s" for reading: %s`, src, err)
		return
	}
	defer fd.Close()

	// One of the possible optimisations here is to split file in chunks and process
	// each chunk individually.
	words, err := collectWords(fd)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `extract: reading input "%s": %s`, src, err)
		return
	}

	if sortResults {
		less := collate.IndexString(lang)
		sort.Slice(words, func(i, j int) bool {
			return less(words[i], words[j])
		})
	}

	wd, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePerm)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `extract: opening destination file "%s" for writing: %s`, src, err)
		return
	}
	defer fd.Close()

	// Writing word by word can result in too many writes, hence, it is slow.
	// Let's add some steroids ...
	wr := bufio.NewWriter(wd)

	if err = writeResults(wr, words); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `extract: writing results "%s": %s`, dst, err)
		return
	}
	if err = wr.Flush(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `extract: writing results "%s": %s`, dst, err)
		return
	}

	// _, _ = fmt.Fprintf(os.Stdout, "Saved %s\n", dst)
}
func collectWords(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)

	scanner.Split(splitWordsUnicode)

	dict := make(map[string]empty)
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		if _, ok := dict[word]; ok {
			continue // duplicate detected
		}

		dict[word] = empty{}

		// Theoretically, if sorting is not needed, we can write right here and
		// skip words slice preparation below.
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// This is expensive ...
	words := make([]string, len(dict))

	var i int
	for w := range dict {
		words[i] = w
		i++

		delete(dict, w)
	}

	return words, nil
}

func writeResults(w io.Writer, words []string) error {
	// This is to preallocate memory once for "string => []byte + \n" conversion
	// and reuse it on every iteration.
	var buf bytes.Buffer
	for _, word := range words {
		buf.WriteString(word)
		buf.WriteRune('\n')

		if _, err := buf.WriteTo(w); err != nil {
			return err
		}

		buf.Reset()
	}

	return nil
}

func ExtractUniqueWords(content string, lang string) ([]string, error) {
	r := strings.NewReader(content)
	words, err := collectWords(r)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `collectWords error: %s`, err)
		return nil, err
	}
	less := collate.IndexString(lang)
	sort.Slice(words, func(i, j int) bool {
		return less(words[i], words[j])
	})
	return words, nil
}
