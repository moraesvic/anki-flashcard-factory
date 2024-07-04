package input

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func GetLines(filepath string) <-chan string {
	ch := make(chan string)
	fp, err := os.Open(filepath)

	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	go func() {
		defer close(ch)
		defer fp.Close()

		scanner := bufio.NewScanner(fp)

		for {
			result := scanner.Scan()
			if !result {
				err := scanner.Err()
				if err != nil {
					log.Println("An error occurred scanning the file.")
					log.Println(err)
				}
				break
			}

			line := scanner.Text()
			lineTrimmed := strings.TrimSpace(line)

			if len(lineTrimmed) == 0 {
				log.Println("Skipping empty line.")
				continue
			}

			ch <- lineTrimmed
		}
	}()

	return ch
}
