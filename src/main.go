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
	TEMPO = 0.7
)

func _toPinyin(s string) string {
	pinyinArgs := pinyin.NewArgs()
	pinyinArgs.Style = pinyin.Tone

	result := pinyin.LazyPinyin(s, pinyinArgs)
	return strings.Join(result, " ")
}

func toPinyin(ch chan Sentence, sentence Sentence) {
	sentence.textTransliterated = _toPinyin(sentence.textOriginal)
	ch <- sentence
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

func changeAudioTempo(ch chan Sentence, sentence Sentence) {
	sentence.audioReducedSpeed = _changeAudioTempo(sentence.audioOriginal, TEMPO)
	ch <- sentence
}

func synthesizeSpeech(ch chan Sentence, sentence Sentence, client *polly.Client) {
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
	ch <- sentence
}

type Sentence struct {
	id                 string
	textOriginal       string
	textTransliterated string
	audioOriginal      string
	audioReducedSpeed  string
}

func printSentence(sentence Sentence) {
	fmt.Printf("%-20s : %-20s\n", "ID", sentence.id)
	fmt.Printf("%-20s : %-20s\n", "Original text", sentence.textOriginal)
	fmt.Printf("%-20s : %-20s\n", "Transliterated text", sentence.textTransliterated)
	fmt.Printf("%-20s : %-20s\n", "Original audio", sentence.audioOriginal)
	fmt.Printf("%-20s : %-20s\n", "Reduced speed audio", sentence.audioReducedSpeed)
	fmt.Println()
}

func makeId(timestamp string, index int) string {
	return fmt.Sprintf("%s-%04d", timestamp, index)
}

func makeSentence(ch chan Sentence, timestamp string, index int, textOriginal string) {
	sentence := Sentence{}
	sentence.id = makeId(timestamp, index)
	sentence.textOriginal = textOriginal
	ch <- sentence
}

func getLines() []string {
	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return lines
}

func getClient() *polly.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Cannot load AWS config!")
	}

	client := polly.NewFromConfig(cfg)
	return client
}

func main() {
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Printf("Starting program at %s\n", timestamp)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file_path>")
		os.Exit(1)
	}

	client := getClient()
	lines := getLines()

	pipe1 := make(chan Sentence, len(lines))
	pipe2 := make(chan Sentence, len(lines))
	pipe3 := make(chan Sentence, len(lines))
	pipe4 := make(chan Sentence, len(lines))

	for index, text := range lines {
		go makeSentence(pipe1, timestamp, index, text)
	}

	for range lines {
		go synthesizeSpeech(pipe2, <-pipe1, client)
	}

	for range lines {
		go changeAudioTempo(pipe3, <-pipe2)
	}

	for range lines {
		go toPinyin(pipe4, <-pipe3)
	}

	for range lines {
		printSentence(<-pipe4)
	}
}
