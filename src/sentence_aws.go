// An implementation for ISentence using AWS Translate, AWS Polly, github.com/mozillazg/go-pinyin, and github.com/moraesvic/shenme

package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"github.com/moraesvic/flashcard-factory/aws"
	"github.com/moraesvic/flashcard-factory/pinyin"
	"github.com/moraesvic/flashcard-factory/types"

	shenme "github.com/moraesvic/shenme/lib"
)

type SentenceAWS struct {
	Sentence
}

var pollyClient *polly.Client
var translateClient *translate.Client

func init() {
	pollyClient = aws.GetPollyClient()
	translateClient = aws.GetTranslateClient()
}

func (SentenceAWS) SynthesizeSpeech(id string, text string) string {
	bytes := aws.SynthesizeSpeech(pollyClient, text)
	audioFile := fmt.Sprintf("%s.mp3", id)

	err := os.WriteFile(audioFile, bytes, 0644)

	if err != nil {
		panic(fmt.Sprintf("Error writing audio data to file: %s", err))
	}

	return audioFile
}

func (SentenceAWS) Translate(text string) string {
	return aws.Translate(translateClient, text)
}

func (SentenceAWS) Pinyin(text string) string {
	return pinyin.Pinyin(text)
}

func (SentenceAWS) ChangeAudioTempo(audioFile string) string {
	return aws.ChangeAudioTempo(audioFile)
}

func (SentenceAWS) DefineHTML(traditional string) types.IDefinitionHTML {
	return shenme.Definer{}.DefineHTML(traditional)
}

func (SentenceAWS) Traditional(simplified string) string {
	return shenme.Traditional(simplified)
}

func (SentenceAWS) WikiURL(traditional string) string {
	return shenme.WikiURL(traditional)
}
