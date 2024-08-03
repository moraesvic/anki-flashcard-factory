package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/moraesvic/flashcard-factory/input"
)

const (
	/*
	 Avoid extreme parallelism, AWS services can be very sensitive to this.

	 https://docs.aws.amazon.com/translate/latest/dg/what-is-limits.html#limits
	 https://docs.aws.amazon.com/polly/latest/dg/limits.html#limits-throttle
	*/
	N_WORKERS = 5
	/*
	 In the future we may want to use other back ends that could theoretically
	 provided better audio synthesis, translation etc.
	*/
	BACKEND = "aws"
)

var start int64
var timestamp string

func init() {
	start = time.Now().UnixMilli()
	timestamp = time.UnixMilli(start).UTC().Format(time.RFC3339)
}

func main() {
	log.Printf("Starting program with up to %d workers...", N_WORKERS)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file_path>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	log.Printf("Reading from file \"%s\"", inputFile)

	lines := input.GetLines(inputFile)

	var wg sync.WaitGroup

	var counter atomic.Uint32
	for range N_WORKERS {
		wg.Add(1)
	}

	fmt.Println(`#separator:tab
#html:true
#tags column:8`)

	for range N_WORKERS {
		go func() {
			defer wg.Done()

			for {
				line, ok := <-lines
				if !ok {
					return
				}

				index := counter.Add(1)
				sentence := Sentence{}.New(timestamp, index, line)
				flashcard := sentence.Flashcard(BACKEND)
				fmt.Println(AnkiString(flashcard))
			}
		}()
	}

	wg.Wait()
	finalCount := counter.Load()

	end := time.Now().UnixMilli()
	ellapsedSeconds := (float64(end - start)) / 1000.0
	averageProcessingTime := float64(finalCount) / ellapsedSeconds

	log.Printf("Processed %d flashcards in %.2f seconds (%.2f cards/s)", finalCount, ellapsedSeconds, averageProcessingTime)
}
