package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	shenme "github.com/moraesvic/shenme/types"
)

type Sentence struct {
	id   string
	text string
}

type PrintableSentence struct {
	id                    string
	text                  string
	transliteration       string
	translation           string
	audioFile             string
	reducedSpeedAudioFile string
}

type ISynthesizeSpeech interface {
	SynthesizeSpeech(id string, text string) (audioFile string)
}

type ITranslate interface {
	Translate(text string) (translation string)
}

type IPinyin interface {
	Pinyin(text string) (transliteration string)
}

type IChangeAudioTempo interface {
	ChangeAudioTempo(audioFile string) (reducedSpeedAudioFile string)
}

type IDefine interface {
	shenme.IDefinerHTML
}

type Flashcard interface {
	ISynthesizeSpeech
	IChangeAudioTempo
	ITranslate
	IPinyin
	IDefine
	Id() string
	Text() string
}

func (s Sentence) Id() string {
	return s.id
}

func (s Sentence) Text() string {
	return s.text
}

func (Sentence) New(timestamp string, index uint32, text string) Sentence {
	id := fmt.Sprintf("%s-%04d", timestamp, index)

	sentence := Sentence{
		id:   id,
		text: text,
	}

	return sentence
}

func (s Sentence) Flashcard(backend string) Flashcard {
	switch backend {
	case "aws":
		return SentenceAWS{s}
	default:
		panic(fmt.Sprintf("Unknown backend %s", backend))
	}
}

func (ps PrintableSentence) String() string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "ID", ps.id))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Text", ps.text))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Transliteration", ps.transliteration))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Translation", ps.translation))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Original audio", ps.audioFile))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Reduced speed audio", ps.reducedSpeedAudioFile))

	return sb.String()
}

func AnkiString(s Flashcard) string {
	id := s.Id()
	text := s.Text()

	var transliteration string
	var audioFile string
	var reducedSpeedAudioFile string
	var translation string

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		transliteration = s.Pinyin(text)
	}()

	go func() {
		defer wg.Done()
		translation = s.Translate(text)
		if len(translation) == 0 {
			translation = "(add translation here)"
		}
	}()

	go func() {
		defer wg.Done()
		audioFile = s.SynthesizeSpeech(id, text)
		reducedSpeedAudioFile = s.ChangeAudioTempo(audioFile)
	}()

	wg.Wait()

	printable := PrintableSentence{
		id:                    id,
		text:                  text,
		transliteration:       transliteration,
		translation:           translation,
		audioFile:             audioFile,
		reducedSpeedAudioFile: reducedSpeedAudioFile,
	}

	log.Println(printable)

	return fmt.Sprintf(
		"%s;%s;%s;%s;%s",
		text,
		transliteration,
		fmt.Sprintf("[sound:%s]", audioFile),
		fmt.Sprintf("[sound:%s]", reducedSpeedAudioFile),
		translation,
	)
}
