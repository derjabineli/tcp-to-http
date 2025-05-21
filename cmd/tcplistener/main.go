package main

import (
	"fmt"
	"log"
	"net"

	"github.com/derjabineli/httpfromtcp/internal/request"
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
		request, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error occured. %v\n", err)
			continue
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %v\n", request.RequestLine.Method)
		fmt.Printf("- Target: %v\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v\n", request.RequestLine.HttpVersion)

    fmt.Println("Headers:")
    for key, val := range request.Headers {
      fmt.Printf("- %v: %v\n", key, val)
    }
		fmt.Println("Connection closed")
	}
}
