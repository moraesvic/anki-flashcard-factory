package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/moraesvic/anki-flashcard-factory/audio"
	"github.com/moraesvic/anki-flashcard-factory/input"
)

const (
	// avoid extreme parallelism
	CHANNEL_CAPACITY = 8
)

func main() {
	start := time.Now().UnixMilli()
	timestamp := fmt.Sprint(start)

	log.Printf("Starting program at %s\n", timestamp)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file_path>")
		os.Exit(1)
	}

	pollyClient := audio.GetPollyClient()
	translateClient := audio.GetTranslateClient()

	inputChannel := make(chan string, CHANNEL_CAPACITY)
	outputChannel := make(chan Sentence, CHANNEL_CAPACITY)

	go input.GetLines(os.Args[1], inputChannel)

	index := 0
	for text := range inputChannel {
		go func(text string, index int) {
			sentence := CreateSentence(timestamp, index, text)
			sentence.SynthesizeSpeech(pollyClient)
			sentence.ChangeAudioTempo()
			sentence.ToPinyin()
			sentence.Translate(translateClient)
			sentence.ToAnkiFlashcard()
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

	log.Printf("Processed %d flashcards in %.2f seconds", index, ellapsedSeconds)
}
