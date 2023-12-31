package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
	"github.com/mozillazg/go-pinyin"
)

const (
	// avoid extreme parallelism?
	CHANNEL_CAPACITY = 10
	TEMPO            = 0.7
)

func _toPinyin(s string) string {
	pinyinArgs := pinyin.NewArgs()
	pinyinArgs.Style = pinyin.Tone

	result := pinyin.LazyPinyin(s, pinyinArgs)
	return strings.Join(result, " ")
}

func toPinyin(sentence *Sentence) {
	sentence.textTransliterated = _toPinyin(sentence.textOriginal)
}

func _changeAudioTempo(path string, tempo float64) string {
	if !strings.HasSuffix(path, ".mp3") {
		panic("Only files of MP3 type are allowed!")
	}

	executable := "ffmpeg"
	inputFile := fmt.Sprintf("file:%s", path)
	atempo := fmt.Sprintf("atempo=%.2f", tempo)

	// For quality, lower is better
	_quality := 2
	quality := fmt.Sprint(_quality)

	nameWithoutExtension := strings.Split(path, ".mp3")[0]
	outputPath := fmt.Sprintf("%s_%s.mp3", nameWithoutExtension, atempo)

	cmd := exec.Command(
		executable,
		"-i",
		inputFile,
		"-filter:a",
		atempo,
		"-q:a",
		quality,
		fmt.Sprintf("file:%s", outputPath))

	res, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error processing %s into %s", path, outputPath)
		fmt.Print(string(res))
		fmt.Print()
		log.Fatal(err)
	}

	return outputPath
}

func changeAudioTempo(sentence *Sentence) {
	sentence.audioReducedSpeed = _changeAudioTempo(sentence.audioOriginal, TEMPO)
}

func synthesizeSpeech(sentence *Sentence, client *polly.Client) {
	params := &polly.SynthesizeSpeechInput{
		Engine:       types.EngineNeural,
		OutputFormat: types.OutputFormatMp3,
		Text:         &sentence.textOriginal,
		VoiceId:      types.VoiceIdZhiyu,
	}

	res, err := client.SynthesizeSpeech(context.TODO(), params)
	if err != nil {
		panic("Could not synthesize speech!")
	}

	audio, err := io.ReadAll(res.AudioStream)
	if err != nil {
		panic("Could not read audio stream!")
	}

	audioOriginal := fmt.Sprintf("%s.mp3", sentence.id)

	err = os.WriteFile(audioOriginal, audio, 0644)
	if err != nil {
		panic(fmt.Sprintf("Error writing audio data to file: %s", err))
	}

	sentence.audioOriginal = audioOriginal
}

type Sentence struct {
	id                 string
	textOriginal       string
	textTransliterated string
	audioOriginal      string
	audioReducedSpeed  string
}

func logSentence(sentence Sentence) {
	log.Printf("%-20s : %-20s\n", "ID", sentence.id)
	log.Printf("%-20s : %-20s\n", "Original text", sentence.textOriginal)
	log.Printf("%-20s : %-20s\n", "Transliterated text", sentence.textTransliterated)
	log.Printf("%-20s : %-20s\n", "Original audio", sentence.audioOriginal)
	log.Printf("%-20s : %-20s\n", "Reduced speed audio", sentence.audioReducedSpeed)
	log.Println()
}

func makeId(timestamp string, index int) string {
	return fmt.Sprintf("%s-%04d", timestamp, index)
}

func makeSentence(timestamp string, index int, textOriginal string) Sentence {
	sentence := Sentence{}
	sentence.id = makeId(timestamp, index)
	sentence.textOriginal = textOriginal

	return sentence
}

func getLines(ch chan string) {
	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		ch <- line
	}

	close(ch)
}

func getClient() *polly.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Cannot load AWS config!")
	}

	client := polly.NewFromConfig(cfg)
	return client
}

func toAnki(sentence Sentence) string {
	return fmt.Sprintf(
		"%s;%s;%s;%s;(add translation here)",
		sentence.textOriginal,
		sentence.textTransliterated,
		fmt.Sprintf("[sound:%s]", sentence.audioOriginal),
		fmt.Sprintf("[sound:%s]", sentence.audioReducedSpeed))
}

func main() {
	timestamp := fmt.Sprint(time.Now().Unix())
	log.Printf("Starting program at %s\n", timestamp)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file_path>")
		os.Exit(1)
	}

	client := getClient()

	input := make(chan string, CHANNEL_CAPACITY)
	output := make(chan Sentence, CHANNEL_CAPACITY)

	go getLines(input)

	index := 0
	for text := range input {
		go func(text string, index int) {
			sentence := makeSentence(timestamp, index, text)
			synthesizeSpeech(&sentence, client)
			changeAudioTempo(&sentence)
			toPinyin(&sentence)
			output <- sentence
		}(text, index)

		index++
	}

	for i := 0; i < index; i++ {
		sentence := <-output
		logSentence(sentence)
		fmt.Println(toAnki(sentence))
	}
}
