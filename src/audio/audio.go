package audio

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	pollyTypes "github.com/aws/aws-sdk-go-v2/service/polly/types"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	translateTypes "github.com/aws/aws-sdk-go-v2/service/translate/types"
)

const (
	// For ffmpeg audio quality, lower is better
	QUALITY           float64 = 2
	TEMPO                     = 0.7
	FFMPEG_EXECUTABLE         = "ffmpeg"
	// https://docs.aws.amazon.com/polly/latest/dg/limits.html
	POLLY_THROTTLING_TIMEOUT_SECONDS = 2.0
)

func changeAudioTempo(inputFile string, outputFile string, atempo float64, quality float64) string {
	cmd := exec.Command(
		FFMPEG_EXECUTABLE,
		"-i",
		fmt.Sprintf("file:%s", inputFile),
		"-filter:a",
		fmt.Sprintf("atempo=%.2f", atempo),
		"-q:a",
		fmt.Sprintf("%.2f", quality),
		fmt.Sprintf("file:%s", outputFile))

	res, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error processing %s into %s", inputFile, outputFile)
		fmt.Print(string(res))
		fmt.Print()
		log.Fatal(err)
	}

	return outputFile
}

func ChangeAudioTempo(inputFile string) string {
	if !strings.HasSuffix(inputFile, ".mp3") {
		panic(fmt.Sprintf("Only files of MP3 type are allowed! Received \"%s\"", inputFile))
	}

	nameWithoutExtension := strings.Split(inputFile, ".mp3")[0]
	outputFile := fmt.Sprintf("%s_atempo=%.2f.mp3", nameWithoutExtension, TEMPO)

	return changeAudioTempo(inputFile, outputFile, TEMPO, QUALITY)
}

func SynthesizeSpeech(client *polly.Client, input string) []byte {
	params := &polly.SynthesizeSpeechInput{
		Engine:       pollyTypes.EngineNeural,
		OutputFormat: pollyTypes.OutputFormatMp3,
		Text:         &input,
		VoiceId:      pollyTypes.VoiceIdZhiyu,
	}

	res, err := client.SynthesizeSpeech(context.TODO(), params)
	if err != nil {
		if strings.Contains(err.Error(), "ThrottlingException") {
			log.Printf("AWS returned ThrottlingException, sleeping for %.1f seconds", POLLY_THROTTLING_TIMEOUT_SECONDS)
			time.Sleep(time.Second * POLLY_THROTTLING_TIMEOUT_SECONDS)
			return SynthesizeSpeech(client, input)
		}

		log.Fatal("Could not synthesize speech!\n", err)
	}

	audio, err := io.ReadAll(res.AudioStream)
	if err != nil {
		panic("Could not read audio stream!")
	}

	return audio
}

func GetPollyClient() *polly.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Cannot load AWS config!")
	}

	client := polly.NewFromConfig(cfg)
	return client
}

func GetTranslateClient() *translate.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Cannot load AWS config!")
	}

	client := translate.NewFromConfig(cfg)
	return client
}

func Translate(client *translate.Client, text string) string {
	contentType := "text/plain"

	document := translateTypes.Document{
		Content:     []byte(text),
		ContentType: &contentType,
	}

	// Note: for Chinese Traditional characters use "zh-TW"
	// https://docs.aws.amazon.com/translate/latest/dg/what-is-languages.html
	sourceLanguageCode := "zh"

	targetLanguageCode := "en"

	// Here you can add settings for Brevity, Formality and Profanity.
	// We don't have any preference, so we will just leave it as is.
	settings := translateTypes.TranslationSettings{}

	input := translate.TranslateDocumentInput{
		Document: &document,
		// Note: for Chinese Traditional characters use zh-TW
		SourceLanguageCode: &sourceLanguageCode,
		TargetLanguageCode: &targetLanguageCode,
		Settings:           &settings,
	}

	translatedDocument, err := client.TranslateDocument(context.TODO(), &input)

	if err != nil {
		log.Fatal("An error occurred translating the document: ", err)
	}

	return string(translatedDocument.TranslatedDocument.Content)
}
