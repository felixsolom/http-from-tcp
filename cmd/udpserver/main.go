package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	UDPAddress, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Couldn't resolve UDP address: %v", err)
	}

	UDPConn, err := net.DialUDP("udp", nil, UDPAddress)
	if err != nil {
		log.Fatalf("Coudn't open UDP connection: %v", err)
	}
	defer UDPConn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Unable to read line: %v", err)
			continue
		}
		if err == io.EOF {
			log.Printf("No more lines to read: %v", err)
			break
		}

		_, err = UDPConn.Write([]byte(line))
		if err != nil {
			log.Printf("Couldn't write line to UDP connection: %v", err)
			continue
		}
	}
}
