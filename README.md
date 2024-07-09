# Introduction

[Anki](https://apps.ankiweb.net/) is an excellent piece of software that helps you memorize any subject by use of SRS (Space Repetition System). I use it for learning languages and much more. Give it a try!

When creating flashcards for studying Chinese in Anki, I wanted to add pinyin and audio along the original text and the translation. In addition to that, I wanted to create a slowed down version of the audio, since my listening skills are still poor. It turned out to be a good opportunity to get a taste of the Go language, also leveraging AWS Polly as a speech generator, and AWS Translate as a generator of stub translations.

# How to use this

First, you must have AWS credentials that allow you to invoke Polly's `SynthesizeSpeech` API, as well as AWS Translate's `TranslateDocument` API. Enable these credentials using `aws configure` before running the program.

Second, the program can be run with `go run main.go <FILENAME>` or by generating a binary with `go build main.go` and then running `./flashcard-factory <FILENAME>`. The binary `flashcard-factory` can be copied into a folder in the user's `PATH` for ease of use.

## Expected input

The program takes exactly one command line option, which is the path to a text file containing sentences in Mandarin Chinese. The name or extension of this file is not relevant. Here is one example, let's say this is `input.txt`:

```
我说中文。
她有两只猫。
我住在巴西。
这种药的作用是什么？
觉
```

Every line must contain a sentence in Mandarin Chinese. Empty lines will be ignored in this file. Comments or any other syntax is not supported.

Invoke the program:

```bash
flashcard-factory input.txt
```

This is the expected output (processing occurs asynchronously and actual order of the output may differ):

```
2024/01/28 20:44:08 Starting program...
2024/01/28 20:44:08 Reading from file "sentences"
2024/01/28 20:44:10
ID                   : 1706485448439-0004
Original text        : 觉
Transliterated text  : jué, jiào
Translated text      : Kyaw
Original audio       : 1706485448439-0004.mp3
Reduced speed audio  : 1706485448439-0004_atempo=0.70.mp3

觉;jué, jiào;[sound:1706485448439-0004.mp3];[sound:1706485448439-0004_atempo=0.70.mp3];Kyaw
2024/01/28 20:44:10
ID                   : 1706485448439-0002
Original text        : 我住在巴西。
Transliterated text  : wǒ zhù zài bā xī .
Translated text      : I live in Brazil.
Original audio       : 1706485448439-0002.mp3
Reduced speed audio  : 1706485448439-0002_atempo=0.70.mp3

我住在巴西。;wǒ zhù zài bā xī .;[sound:1706485448439-0002.mp3];[sound:1706485448439-0002_atempo=0.70.mp3];I live in Brazil.
2024/01/28 20:44:10
ID                   : 1706485448439-0000
Original text        : 我说中文。
Transliterated text  : wǒ shuō zhōng wén .
Translated text      : I speak Chinese.
Original audio       : 1706485448439-0000.mp3
Reduced speed audio  : 1706485448439-0000_atempo=0.70.mp3

我说中文。;wǒ shuō zhōng wén .;[sound:1706485448439-0000.mp3];[sound:1706485448439-0000_atempo=0.70.mp3];I speak Chinese.
2024/01/28 20:44:10
ID                   : 1706485448439-0003
Original text        : 这种药的作用是什么？
Transliterated text  : zhè zhǒng yào de zuò yòng shì shén me ?
Translated text      : What is the effect of this medicine?
Original audio       : 1706485448439-0003.mp3
Reduced speed audio  : 1706485448439-0003_atempo=0.70.mp3

这种药的作用是什么？;zhè zhǒng yào de zuò yòng shì shén me ?;[sound:1706485448439-0003.mp3];[sound:1706485448439-0003_atempo=0.70.mp3];What is the effect of this medicine?
2024/01/28 20:44:10
ID                   : 1706485448439-0001
Original text        : 她有两只猫。
Transliterated text  : tā yǒu liǎng zhǐ māo .
Translated text      : She has two cats.
Original audio       : 1706485448439-0001.mp3
Reduced speed audio  : 1706485448439-0001_atempo=0.70.mp3

她有两只猫。;tā yǒu liǎng zhǐ māo .;[sound:1706485448439-0001.mp3];[sound:1706485448439-0001_atempo=0.70.mp3];She has two cats.
2024/01/28 20:44:10 Processed 5 flashcards in 1.75 seconds (2.85 cards/s)
```

The lines starting with a timestamp are for diagnostic purposes only. Usually, we are not interested in them. Since these are printed to stdout, you can get rid of them by invoking the program like this:

```bash
flashcard-factory input.txt 2> /dev/null
```

Thus you will only get:

```
觉;jué, jiào;[sound:1706485481408-0004.mp3];[sound:1706485481408-0004_atempo=0.70.mp3];Kyaw
她有两只猫。;tā yǒu liǎng zhǐ māo .;[sound:1706485481408-0001.mp3];[sound:1706485481408-0001_atempo=0.70.mp3];She has two cats.
我说中文。;wǒ shuō zhōng wén .;[sound:1706485481408-0000.mp3];[sound:1706485481408-0000_atempo=0.70.mp3];I speak Chinese.
这种药的作用是什么？;zhè zhǒng yào de zuò yòng shì shén me ?;[sound:1706485481408-0003.mp3];[sound:1706485481408-0003_atempo=0.70.mp3];What is the effect of this medicine?
我住在巴西。;wǒ zhù zài bā xī .;[sound:1706485481408-0002.mp3];[sound:1706485481408-0002_atempo=0.70.mp3];I live in Brazil.
```

You will quickly notice that this is a CSV format using a semicolon as a delimiter. The columns correspond to a note type I created in Anki, receiving the original text, pinyin transcription, original audio, slowed down audio, and translation.

**Please note (1)**: If a line of input consists of only one Chinese character, the program will yield all the possible pinyin readings. If a line contains punctuation, the punctuation will be preserved in the pinyin output.

**Please note (2)**: The generated pinyin and translation are the best guess that the tools can provide given an input with very little context. Please make sure to review them yourself before you add the flashcards to your library.

By copying and pasting the CSV lines to a file and saving it, they can now be imported into Anki by going to `File > Import`.

You will also notice that media (MP3) files were generated:

```
viih@viih-samsung:~/flashcard-factory/src$ ls -1
'1706485481408-0000_atempo=0.70.mp3'
1706485481408-0000.mp3
'1706485481408-0001_atempo=0.70.mp3'
1706485481408-0001.mp3
'1706485481408-0002_atempo=0.70.mp3'
1706485481408-0002.mp3
'1706485481408-0003_atempo=0.70.mp3'
1706485481408-0003.mp3
'1706485481408-0004_atempo=0.70.mp3'
1706485481408-0004.mp3
```

These must be copied into your Anki's media folder. Please consult the Anki manual in case of questions.

## End result

When all of this is done, you can see (and hear!) your new flashcards in Anki:

![Illustration of the final result (1)](/demo-1.png)

![Illustration of the final result (2)](/demo-2.png)

# Testing the program

The attached file `test_data.txt` contains 198 sentences in Chinese that can be used to test the correctness and performance of the program.

# Task-level parallelization

Testing with a large enough data set, task-level parallelization was able to improve processing times in about 17%:

```
# With task-level parallelization
2024/07/09 13:59:16 Processed 198 flashcards in 22.35 seconds (8.86 cards/s)
# Without task-level parallelization
2024/07/09 14:01:26 Processed 198 flashcards in 26.21 seconds (7.56 cards/s)
```

# Docker

It is possible, though not very convenient, to run the program from Docker. Build the Docker image and tag it (for example, with `flashcard-factory`). This is usually not worth the trouble, because Go can already generate a standalone executable, but may be a possibility if you don't have FFMPEG installed in your machine or have other compatibility problems.

Then run a container with that image, mounting the current directory to the Docker working directory:

```
docker build -t flashcard-factory -f docker/Dockerfile .
docker run -v $PWD:$PWD -w $PWD flashcard-factory
```

Then, inside Docker you will need to add your AWS credentials.

Finally, you can invoke `flashcard-factory` as explained above.

# Caveats

One single Chinese character can have more than one reading depending on the context (for example, 了). The library we are using to do the pinyin transcription can offer all of the possibilities for transcription, but cannot pick the right one. For the sake of simplicity, we pick the most likely transcription, which can be easily edited in the import phase.

## Punctuation

This kind of output looks weird and should be fixed.

```
Original text        : 我不说英文，我只说汉语。
Transliterated text  : wǒ bù shuō yīng wén ,wǒ zhǐ shuō hàn yǔ .
```

Punctuation should come right after the word, as usual:

# Note on SoX vs. FFMPEG

SoX is not recommended, sounds bad even when you turn on the best quality:

```bash
sox polly.mp3 -C 48.0 sox.mp3 tempo 0.75
```

# TODO

## Testing

We still need to add tests. Most of the code is glue code calling external libraries and APIs, but there is a little bit of business logic that we would like to assure is doing the right thing.

## New features

- **Conversion to/from Traditional Chinese characters**: https://github.com/siongui/gojianfan
- **Importing other formats than the plain text file, for example, CSV**: https://pkg.go.dev/encoding/csv
- The feature above could be used to work with datasets from Tatoeba, for example
- **Reading job specifications from a file, or from command line args**: https://pkg.go.dev/flag
