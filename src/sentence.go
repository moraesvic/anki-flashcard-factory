package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/moraesvic/flashcard-factory/types"
)

type Sentence struct {
	id   string
	text string
}

type PrintableSentence struct {
	id          string
	text        string
	traditional string
	// `kind` is either `sentence` or `vocabulary`
	kind            string
	transliteration string
	// only defined when `kind` = `sentence`
	translation string
	// only defined when `kind` = `vocabulary`
	definition string
	// only defined when `kind` = `vocabulary`
	wikiURL               string
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
	types.IDefinerHTML
}

type ITraditional interface {
	Traditional(simplified string) string
}

type IWikiURL interface {
	WikiURL(traditional string) string
}

type Flashcard interface {
	ISynthesizeSpeech
	IChangeAudioTempo
	ITranslate
	IPinyin
	IDefine
	ITraditional
	IWikiURL
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

func (ps PrintableSentence) output() string {
	if ps.kind == "sentence" {
		return ps.translation
	} else {
		return ps.definition
	}
}

func (ps PrintableSentence) String() string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "ID", ps.id))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Text", ps.text))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Traditional", ps.traditional))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Kind", ps.kind))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Transliteration", ps.transliteration))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Translation", ps.translation))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Definition", ps.definition))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Output", ps.output()))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Wiki URL", ps.wikiURL))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Original audio", ps.audioFile))
	sb.WriteString(fmt.Sprintf("%-20s : %-20s\n", "Reduced speed audio", ps.reducedSpeedAudioFile))

	return sb.String()
}

func (ps PrintableSentence) Anki() string {
	return fmt.Sprintf(
		"%s\t%s\t%s\t%s\t%s",
		ps.text,
		ps.transliteration,
		fmt.Sprintf("[sound:%s]", ps.audioFile),
		fmt.Sprintf("[sound:%s]", ps.reducedSpeedAudioFile),
		ps.output(),
	)
}

const CJKPunctuation = "？！，、。（）：【】;"

// We will consider Kind=`sentence` when the text is more than 4 characters long,
// or when the string contains any kind of punctuation. Otherwise, Kind=`vocabulary`.
func Kind(f Flashcard) string {
	if strings.ContainsAny(f.Text(), CJKPunctuation) || len([]rune(f.Text())) > 4 {
		return "sentence"
	} else {
		return "vocabulary"
	}
}

func AnkiString(f Flashcard) string {
	id := f.Id()
	text := f.Text()
	kind := Kind(f)
	traditional := f.Traditional(text)

	var transliteration string
	var audioFile string
	var reducedSpeedAudioFile string
	var translation string
	var definition string
	var wikiURL string

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		transliteration = f.Pinyin(text)
	}()

	go func() {
		defer wg.Done()

		if kind == "sentence" {
			translation = f.Translate(text)
			if len(translation) == 0 {
				translation = "(add translation here)"
			}
		} else {
			wikiURL = f.WikiURL(traditional)
			definitions := f.DefineHTML(traditional)

			if definitions.Length() == 0 {
				definition = "(add definition here)"
			} else {
				definition = definitions.HTML()
			}
		}
	}()

	go func() {
		defer wg.Done()
		audioFile = f.SynthesizeSpeech(id, text)
		reducedSpeedAudioFile = f.ChangeAudioTempo(audioFile)
	}()

	wg.Wait()

	printable := PrintableSentence{
		id:                    id,
		text:                  text,
		traditional:           traditional,
		kind:                  kind,
		transliteration:       transliteration,
		translation:           translation,
		definition:            definition,
		wikiURL:               wikiURL,
		audioFile:             audioFile,
		reducedSpeedAudioFile: reducedSpeedAudioFile,
	}

	log.Println(printable)

	return printable.Anki()
}
