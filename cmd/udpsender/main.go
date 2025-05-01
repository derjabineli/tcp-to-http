package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const address = "localhost:42069"

func main() {
	udpEndPoint, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, udpEndPoint)
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()

	buf := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print(">")
		string, err := buf.ReadString('\n')
		if err != nil {
			log.Print(err.Error())
		}
		conn.Write([]byte(string))
	}
}