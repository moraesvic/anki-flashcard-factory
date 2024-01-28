# Introduction

[Anki](https://apps.ankiweb.net/) is an excellent piece of software that helps you memorize any subject by use of SRS (Space Repetition System). I use it for learning languages and much more. Give it a try!

When creating flashcards for studying Chinese in Anki, I wanted to add pinyin and audio along the original text and the translation. In addition to that, I wanted to create a slowed down version of the audio, since my listening skills are still poor. It turned out to be a good opportunity to get a taste of the Go language, also leveraging AWS Polly as a speech generator.

# How to use this

First, you must have AWS credentials that allow you to invoke Polly's `SynthesizeSpeech` API. Enable these credentials using `aws configure` before running the program.

Second, the program can be run with `go run main.go <FILENAME>` or by generating a binary with `go build main.go` and then running `./flashcard-factory <FILENAME>`. The binary `flashcard-factory` can be copied into the user's `bin` folder for ease of use.

## Expected input

The program takes exactly one command line option, which is the path to a text file containing sentences in Mandarin Chinese. The name or extension of this file is not relevant. Here is one example, let's say this is `input.txt`:

```
我说中文。
她友两只猫。
我主在巴西。
这种药的作用是什么？
```

Every line must contain a sentence in Mandarin Chinese. Empty lines will be ignored in this file. Comments or any other syntax is not supported.

Invoke the program:

```bash
flashcard-factory input.txt
```

This is the expected output (processing occurs asynchronously and actual order of the output may differ):

```
2024/01/10 15:47:26 Starting program at 1704912446
2024/01/10 15:47:26 Skipping empty line.
2024/01/10 15:47:27 ID                   : 1704912446-0001
2024/01/10 15:47:27 Original text        : 她友两只猫。
2024/01/10 15:47:27 Transliterated text  : tā yǒu liǎng zhǐ māo
2024/01/10 15:47:27 Original audio       : 1704912446-0001.mp3
2024/01/10 15:47:27 Reduced speed audio  : 1704912446-0001_atempo=0.70.mp3
2024/01/10 15:47:27
她友两只猫。;tā yǒu liǎng zhǐ māo;[sound:1704912446-0001.mp3];[sound:1704912446-0001_atempo=0.70.mp3];(add translation here)
2024/01/10 15:47:27 ID                   : 1704912446-0002
2024/01/10 15:47:27 Original text        : 我主在巴西。
2024/01/10 15:47:27 Transliterated text  : wǒ zhǔ zài bā xī
2024/01/10 15:47:27 Original audio       : 1704912446-0002.mp3
2024/01/10 15:47:27 Reduced speed audio  : 1704912446-0002_atempo=0.70.mp3
2024/01/10 15:47:27
我主在巴西。;wǒ zhǔ zài bā xī;[sound:1704912446-0002.mp3];[sound:1704912446-0002_atempo=0.70.mp3];(add translation here)
2024/01/10 15:47:27 ID                   : 1704912446-0000
2024/01/10 15:47:27 Original text        : 我说中文。
2024/01/10 15:47:27 Transliterated text  : wǒ shuō zhōng wén
2024/01/10 15:47:27 Original audio       : 1704912446-0000.mp3
2024/01/10 15:47:27 Reduced speed audio  : 1704912446-0000_atempo=0.70.mp3
2024/01/10 15:47:27
我说中文。;wǒ shuō zhōng wén;[sound:1704912446-0000.mp3];[sound:1704912446-0000_atempo=0.70.mp3];(add translation here)
2024/01/10 15:47:27 ID                   : 1704912446-0003
2024/01/10 15:47:27 Original text        : 这种药的作用是什么？
2024/01/10 15:47:27 Transliterated text  : zhè zhǒng yào de zuò yòng shì shén me
2024/01/10 15:47:27 Original audio       : 1704912446-0003.mp3
2024/01/10 15:47:27 Reduced speed audio  : 1704912446-0003_atempo=0.70.mp3
2024/01/10 15:47:27
这种药的作用是什么？;zhè zhǒng yào de zuò yòng shì shén me;[sound:1704912446-0003.mp3];[sound:1704912446-0003_atempo=0.70.mp3];(add translation here)
```

The lines starting with a timestamp are for diagnostic purposes only. Usually, we are not interested in them. Since these are printed to stdout, you can get rid of them by invoking the program like this:

```bash
flashcard-factory input.txt 2> /dev/null
```

Thus you will only get:

```
我主在巴西。;wǒ zhǔ zài bā xī;[sound:1704912625-0002.mp3];[sound:1704912625-0002_atempo=0.70.mp3];(add translation here)
我说中文。;wǒ shuō zhōng wén;[sound:1704912625-0000.mp3];[sound:1704912625-0000_atempo=0.70.mp3];(add translation here)
她友两只猫。;tā yǒu liǎng zhǐ māo;[sound:1704912625-0001.mp3];[sound:1704912625-0001_atempo=0.70.mp3];(add translation here)
这种药的作用是什么？;zhè zhǒng yào de zuò yòng shì shén me;[sound:1704912625-0003.mp3];[sound:1704912625-0003_atempo=0.70.mp3];(add translation here)
```

You will quickly notice that this is a CSV format using a semicolon as a delimiter. The columns correspond to a note type I created in Anki, receiving the original text, pinyin transcription, original audio, slowed down audio, and translation.

The translation will have to be added by yourself.

By copying and pasting the CSV lines to a file and saving it, they can now be imported into Anki by going to `File > Import`.

You will also notice that media (MP3) files were generated:

```
viih@viih-samsung:~/anki-flashcard-factory/src$ ls -1
'1704912625-0000_atempo=0.70.mp3'
1704912625-0000.mp3
'1704912625-0001_atempo=0.70.mp3'
1704912625-0001.mp3
'1704912625-0002_atempo=0.70.mp3'
1704912625-0002.mp3
'1704912625-0003_atempo=0.70.mp3'
1704912625-0003.mp3
```

These must be copied into your Anki's media folder. Please consult the Anki manual in case of questions.

## End result

When all of this is done, you can see (and hear!) your new flashcards in Anki:

![Illustration of the final result (1)](/demo-1.png)

![Illustration of the final result (2)](/demo-2.png)

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

# Note on SoX vs. FFMPEG

SoX is not recommended, sounds bad even when you turn on the best quality:

```bash
sox polly.mp3 -C 48.0 sox.mp3 tempo 0.75
```
