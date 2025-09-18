package main

import (
	"fmt"
	"log"
	"net"

	"github.com/felixsolom/http-from-tcp/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Couldn't create listener, %v", err)
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Couldn't accept TCP connection: %v", err)
		}

		fmt.Println("====================================")
		fmt.Println("Accepted connection from:", conn.RemoteAddr())

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("Couldn't get request-line from reader: %v", err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %v\n", request.RequestLine.Method)
		fmt.Printf("- Target: %v\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v\n", request.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key, value := range request.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		fmt.Println("Successfully finished printing lines from connection")
		fmt.Println("Connection to:", conn.RemoteAddr(), "was closed")
	}
}
