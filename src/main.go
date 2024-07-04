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
	// Avoid extreme parallelism, AWS Translate can be very sensitive to this
	N_WORKERS = 5
)

var start int64
var timestamp string

func init() {
	start = time.Now().UnixMilli()
	timestamp = fmt.Sprint(start)
}

func Work(counter *atomic.Uint32, lines <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		line, ok := <-lines
		if !ok {
			return
		}

		index := counter.Add(1)
		sentence := CreateSentence(timestamp, index, line)
		sentence.Process()
		sentence.Log()
		fmt.Println(sentence.ankiFlashcard)
	}
}

func main() {
	log.Println("Starting program...")

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
		go Work(&counter, lines, &wg)
	}

	wg.Wait()
	finalCount := counter.Load()

	end := time.Now().UnixMilli()
	ellapsedSeconds := (float64(end) - float64(start)) / 1000.0
	averageProcessingTime := float64(finalCount) / ellapsedSeconds

	log.Printf("Processed %d flashcards in %.2f seconds (%.2f cards/s)", finalCount, ellapsedSeconds, averageProcessingTime)
}
