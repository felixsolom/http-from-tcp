package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	defer file.Close()

	ch := getLinesChannel(file)
	for line := range ch {
		fmt.Println("read:" + line)
	}

	fmt.Println("Successfully finished reading the file")
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
				log.Printf("End of file reached: %v", err)
				break
			}
			if err != nil {
				log.Printf("unable to read from file: %v", err)
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
