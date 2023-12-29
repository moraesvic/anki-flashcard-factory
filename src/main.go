package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
	"github.com/mozillazg/go-pinyin"
)

func toPinyin(s string) string {
	pinyinArgs := pinyin.NewArgs()
	pinyinArgs.Style = pinyin.Tone

	result := pinyin.LazyPinyin(s, pinyinArgs)
	return strings.Join(result, " ")
}

func changeAudioTempo(path string, tempo float32) string {
	if !strings.HasSuffix(path, ".mp3") {
		panic("Only files of MP3 type are allowed!")
	}

	executable := "ffmpeg"
	inputFile := path
	atempo := fmt.Sprintf("atempo=%f", tempo)
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
		outputPath)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)

	return outputPath
}

func main() {
	fmt.Println("Starting program...")

	chinese := "我说中文。我女朋友叫"
	fmt.Println(toPinyin(chinese))

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Cannot load AWS config!")
	}

	client := polly.NewFromConfig(cfg)
	fmt.Println(client)

	// Create input parameters for the SynthesizeSpeech operation
	params := &polly.SynthesizeSpeechInput{
		Engine:       types.EngineNeural,
		OutputFormat: types.OutputFormatMp3,
		Text:         &chinese,
		VoiceId:      types.VoiceIdZhiyu,
	}

	res, err := client.SynthesizeSpeech(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	// changeAudioTempo("hello.mp3", 0.75)
}
