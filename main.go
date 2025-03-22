package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {

	filename := "messages.txt"
	chunkSize := 8

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, chunkSize)
	curLine := ""

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				if len(curLine) > 0 {
					fmt.Printf("read: %s\n", curLine)
				}
				break
			}
			fmt.Println("Error reading file:", err)
			return
		}
		if n > 0 {
			splitLine := strings.Split(string(buffer[:n]), "\n")
			for i, line := range splitLine {
				if i == len(splitLine)-1 {
					curLine += line
					break
				}
				fmt.Printf("read: %s%s\n", curLine, line)
				curLine = ""
			}
		}
	}
}
