package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const chunkSize = 8
const filePath = "messages.txt"
const listenPort = ":42069"

func main() {

	listner, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("unable to net.listen: %v", err)
	}
	defer listner.Close()

	fmt.Println("Listening on ", listenPort)

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatalf("tcp accept error: %v", err)
		}
		fmt.Printf("Connection accepted %s \n", conn.RemoteAddr())

		go func() {
			defer conn.Close()
			linesChan := getLinesChannel(conn)
			for line := range linesChan {
				fmt.Println(line)
			}
		}()

	}
}

func getLinesChannel(conn io.ReadCloser) <-chan string {
	linesOutput := make(chan string)

	go func() {
		defer conn.Close()
		defer close(linesOutput)
		currentLineContents := ""
		for {
			b := make([]byte, 8, 8)
			n, err := conn.Read(b)
			if err != nil {
				if currentLineContents != "" {
					linesOutput <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			str := string(b[:n])
			parts := strings.Split(str, "\n")
			for i := range len(parts) - 1 {
				linesOutput <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return linesOutput
}
