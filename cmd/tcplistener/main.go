package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Established Connection")
		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Println(line)
		}
		fmt.Println("Connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChannel := make(chan string)
	go func ()  {
		defer f.Close()
		defer close(linesChannel)
		currentLine := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					linesChannel <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				return
			}
	
			section := string(buffer[:n])
			parts := strings.Split(section, "\n")
			if len(parts) == 1 {
				currentLine += parts[0]
				continue
			}
			for i := 0; i < len(parts) - 1; i++ {
				currentLine += parts[i]
			}
			linesChannel <- currentLine
			currentLine = parts[len(parts) - 1]
		}
	}()
	
	return linesChannel
}
