package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"github.com/moraesvic/flashcard-factory/aws"
	"github.com/moraesvic/flashcard-factory/pinyin"
)

type Sentence struct {
	id                    string
	textOriginal          string
	textTransliterated    string
	textTranslated        string
	audioOriginalBytes    []byte
	audioOriginalFile     string
	audioReducedSpeedFile string
	ankiFlashcard         string
}

var pollyClient *polly.Client
var translateClient *translate.Client

func init() {
	pollyClient = aws.GetPollyClient()
	translateClient = aws.GetTranslateClient()
}

func createSentenceId(timestamp string, index uint32) string {
	return fmt.Sprintf("%s-%04d", timestamp, index)
}

func CreateSentence(timestamp string, index uint32, textOriginal string) Sentence {
	sentence := Sentence{
		id:           createSentenceId(timestamp, index),
		textOriginal: textOriginal,
	}

	return sentence
}

func (s Sentence) ToString() string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "ID", s.id))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Original text", s.textOriginal))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Transliterated text", s.textTransliterated))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Translated text", s.textTranslated))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Original audio", s.audioOriginalFile))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Reduced speed audio", s.audioReducedSpeedFile))

	return sb.String()
}

func (s *Sentence) SynthesizeSpeech() {
	s.audioOriginalBytes = aws.SynthesizeSpeech(pollyClient, s.textOriginal)
	s.audioOriginalFile = fmt.Sprintf("%s.mp3", s.id)

	err := os.WriteFile(s.audioOriginalFile, s.audioOriginalBytes, 0644)
	if err != nil {
		panic(fmt.Sprintf("Error writing audio data to file: %s", err))
	}
}

func (s *Sentence) Translate() {
	s.textTranslated = aws.Translate(translateClient, s.textOriginal)
}

func (s Sentence) Log() {
	log.Println(s.ToString())
}

func (s *Sentence) ToAnkiFlashcard() {
	var translation string

	if len(s.textTranslated) == 0 {
		translation = "(add translation here)"
	} else {
		translation = s.textTranslated
	}

	s.ankiFlashcard = fmt.Sprintf(
		"%s;%s;%s;%s;%s",
		s.textOriginal,
		s.textTransliterated,
		fmt.Sprintf("[sound:%s]", s.audioOriginalFile),
		fmt.Sprintf("[sound:%s]", s.audioReducedSpeedFile),
		translation)
}

func (s *Sentence) ToPinyin() {
	s.textTransliterated = pinyin.ToPinyin(s.textOriginal)
}

func (s *Sentence) ChangeAudioTempo() {
	s.audioReducedSpeedFile = aws.ChangeAudioTempo(s.audioOriginalFile)
}

func (s *Sentence) Process() {
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		s.SynthesizeSpeech()
		s.ChangeAudioTempo()
	}()

	go func() {
		defer wg.Done()
		s.Translate()
	}()

	go func() {
		defer wg.Done()
		s.ToPinyin()
	}()

	wg.Wait()
	s.ToAnkiFlashcard()
}
