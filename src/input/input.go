package input

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func ScanLines(fp io.Reader, ch chan string) {
	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		line := scanner.Text()
		lineTrimmed := strings.TrimSpace(line)

		if len(lineTrimmed) == 0 {
			log.Println("Skipping empty line.")
			continue
		}

		ch <- lineTrimmed
	}

	close(ch)
}

func GetLines(file string, ch chan string) {
	fp, err := os.Open(file)
	defer fp.Close()

	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	ScanLines(fp, ch)
}
