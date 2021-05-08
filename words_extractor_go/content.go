package main

import (
	"bufio"
	"bytes"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func getRows(metaPath string) ListOfStrings {
	path := strings.Replace(metaPath, ".yml", ".txt", -1)
	data, _ := os.Open(path)
	defer data.Close()

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	return txtlines
}

func (arr ListOfStrings) toString() string {
	buf := bytes.Buffer{}
	for _, row := range arr {
		text := parseVerse(row).text
		buf.WriteString(text + " ")
	}
	return buf.String()
}

type verseStruct struct {
	book    string
	chapter int
	verse   int
	text    string
}

func parseVerse(s string) verseStruct {
	re := regexp.MustCompile("([^ ]+) ([^ ]+):([^ ]+)\\s?(.*)")
	match := re.FindStringSubmatch(s)
	chapter, chapterError := strconv.Atoi(match[2])
	if chapterError != nil {
		panic(chapterError)
	}
	verse, verseError := strconv.Atoi(match[3])
	if verseError != nil {
		panic(verseError)
	}
	return verseStruct{
		book:    match[1],
		chapter: chapter,
		verse:   verse,
		text:    match[4],
	}
}
