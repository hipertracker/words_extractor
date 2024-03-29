package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"wordextractor/app"
)

const (
	srcPath = "../data/??/**/*.yml"
	outDir  = "words"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot get working directory: %s", err)
		os.Exit(1)
	}

	defaultNumWorkers := runtime.NumCPU()

	var sortResults bool
	var numWorkers int

	flag.IntVar(&numWorkers, "n", defaultNumWorkers, "Number of workers to run (zero to match the number of available CPUs)")
	flag.BoolVar(&sortResults, "s", sortResults, "Sort results")
	flag.Parse()

	if numWorkers <= 0 {
		numWorkers = defaultNumWorkers
	}

	t := time.Now()

	var exitCode int
	if err = app.Run(filepath.Join(wd, srcPath), filepath.Join(wd, outDir), numWorkers, sortResults); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Worker error: %s", err)
		exitCode = 1
	}

	timeTrack(t)
	os.Exit(exitCode)
}

func timeTrack(start time.Time) {
	fmt.Println("Total timing: ", time.Since(start))
}
