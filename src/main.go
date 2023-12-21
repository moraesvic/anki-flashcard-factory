package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		"polly.mp3",
		"-filter:a",
		"atempo=0.75",
		"-q:a",
		// For quality, lower is better
		"2",
		"ffmpeg.mp3")

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
}
