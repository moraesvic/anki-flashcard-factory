package pinyin

import (
	"strings"
	"unicode"

	gopinyin "github.com/mozillazg/go-pinyin"
)

// https://stackoverflow.com/questions/70971932/how-to-check-if-the-rune-is-chinese-punctuation-character-in-go
func convertCJKToWesternPunctuation(s string) string {
	cjkPunctuationToWesternPunctuation := map[rune]rune{
		'？': '?',
		'！': '!',
		'，': ',',
		'、': ',',
		'。': '.',
		'（': '(',
		'）': ')',
		'：': ':',
	}
	output := []rune{}

	for _, r := range s {
		westernPunctuation, ok := cjkPunctuationToWesternPunctuation[r]

		if !ok {
			output = append(output, r)
			continue
		}

		output = append(output, westernPunctuation)
	}

	return string(output)
}

func convertMultiCharacterString(s string) string {
	pinyinArgs := gopinyin.NewArgs()
	pinyinArgs.Style = gopinyin.Tone

	result := gopinyin.LazyPinyin(s, pinyinArgs)
	output := []string{}
	i := 0

	inputWesternized := convertCJKToWesternPunctuation(s)

	for _, r := range inputWesternized {
		if unicode.Is(unicode.Han, r) {
			output = append(output, result[i])
			output = append(output, " ")
			i++
			continue
		}

		output = append(output, string(r))
	}

	return strings.Join(output, "")
}

func convertSingleCharacterString(s string) string {
	pinyinArgs := gopinyin.NewArgs()
	pinyinArgs.Style = gopinyin.Tone
	pinyinArgs.Heteronym = true

	result := gopinyin.Pinyin(s, pinyinArgs)[0]
	return strings.Join(result, ", ")
}

func Pinyin(s string) string {
	if len([]rune(s)) > 1 {
		return convertMultiCharacterString(s)
	} else {
		return convertSingleCharacterString(s)
	}
}
