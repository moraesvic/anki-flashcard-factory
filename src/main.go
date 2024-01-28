package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/moraesvic/anki-flashcard-factory/aws"
	"github.com/moraesvic/anki-flashcard-factory/input"
)

const (
	// Avoid extreme parallelism, AWS Translate can be very sensitive to this
	CHANNEL_CAPACITY = 5
)

func main() {
	start := time.Now().UnixMilli()
	timestamp := fmt.Sprint(start)

	log.Printf("Starting program at %s\n", timestamp)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file_path>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	log.Printf("Reading from file \"%s\"", inputFile)

	pollyClient := aws.GetPollyClient()
	translateClient := aws.GetTranslateClient()

	inputChannel := make(chan string, CHANNEL_CAPACITY)
	outputChannel := make(chan Sentence, CHANNEL_CAPACITY)

	go input.GetLines(inputFile, inputChannel)

	index := 0
	for text := range inputChannel {
		go func(text string, index int) {
			sentence := CreateSentence(timestamp, index, text)
			sentence.Process(pollyClient, translateClient)
			outputChannel <- sentence
		}(text, index)

		index++
	}

	for i := 0; i < index; i++ {
		sentence := <-outputChannel
		sentence.Log()
		fmt.Println(sentence.ankiFlashcard)
	}

	end := time.Now().UnixMilli()
	ellapsedSeconds := (float64(end) - float64(start)) / 1000.0
	averageProcessingTime := float64(index) / ellapsedSeconds

	log.Printf("Processed %d flashcards in %.2f seconds (%.2f cards/s)", index, ellapsedSeconds, averageProcessingTime)
}
