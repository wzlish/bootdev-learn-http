package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const chunkSize = 8
const filePath = "messages.txt"

func main() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}

	linesChan := getLinesChannel(file)
	for line := range linesChan {
		fmt.Println("read:", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesOutput := make(chan string)

	go func() {
		defer f.Close()
		defer close(linesOutput)

		reader := bufio.NewReaderSize(f, chunkSize)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					if line != "" {
						linesOutput <- line
					}
					break
				}
				fmt.Printf("read error: %v\n", err.Error())
				return
			}
			linesOutput <- line
		}
	}()
	return linesOutput
}
