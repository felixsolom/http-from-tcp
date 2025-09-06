package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Printf("Couldn't create listener, %v", err)
	}
	defer listener.Close()
	fmt.Println("Connection was closed")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Couldn't accept TCP connection")
		}

		fmt.Println("====================================")
		fmt.Println("=== Connection Has Been Accepted ===")

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Println(line)
		}
		fmt.Println("Successfully finished printing lines from connection")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		var builder strings.Builder
		buff := make([]byte, 8)

		defer close(ch)
		defer f.Close()

		for {
			n, err := f.Read(buff)
			if err == io.EOF {
				log.Printf("End of lines reached: %v", err)
				break
			}
			if err != nil {
				log.Printf("unable to read from lines: %v", err)
			}

			builder.Write(buff[:n])
		}
		text := builder.String()
		lines := strings.Split(text, "\n")
		for i := 0; i < len(lines); i++ {
			ch <- lines[i]
		}

	}()
	return ch
}
