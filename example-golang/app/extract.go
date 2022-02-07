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
	"unsafe"

	"github.com/tidwall/collate"
)

const (
	filePerm        = 0644
	InitialDictSize = 10000
)

// splitWordsUnicode splits data into words, using Unicode Letter character class.
// It works similar to the regular expression "[^\p{L}]+". This is what was used
// in the original code. Unicode function has slight overhead, but handles UTF-8
// correctly.
//
// Rust and Python versions split text according to "[\W\d]+" - anything that is
// not a word or a digit. TODO: comfirm if some words contain digits
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

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func memHashString(str string) uint64 {
	ss := (*stringStruct)(unsafe.Pointer(&str))
	return uint64(memhash(ss.str, 0, uintptr(ss.len)))
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
	words, err := collectWords(fd, lang, InitialDictSize)
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

	_, _ = fmt.Fprintf(os.Stdout, "Saved %s\n", dst)
}

func collectWords(r io.Reader, lang string, sizeHint int) ([]string, error) {
	scanner := bufio.NewScanner(r)
	ascii := []string{"en", "la"}
	if stringInSlice(lang, ascii) {
		scanner.Split(splitWords)
	} else {
		scanner.Split(splitWordsUnicode)
	}

	// map[uint64]empty should take less memory than map[string]empty and avoid
	// GC checks.
	//
	// sizeHint is used to preallocate map[string]empty and []string slice and skip
	// initial reallocation when they should grow. It is a "magic" number which
	// should not be too big or too small. Ideally, it should be approximated from
	// the text.
	dict := make(map[uint64]empty, sizeHint)
	words := make([]string, 0, sizeHint)

	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		hash := memHashString(word)
		if _, ok := dict[hash]; ok {
			continue // duplicate detected
		}

		dict[hash] = empty{}
		words = append(words, word)

		// Theoretically, if sorting is not needed, we can write right here and
		// skip words slice preparation below.
	}
	if err := scanner.Err(); err != nil {
		return nil, err
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

func ExtractUniqueWords(content string, lang string, sizeHint int) ([]string, error) {
	r := strings.NewReader(content)
	words, err := collectWords(r, lang, sizeHint)

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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
