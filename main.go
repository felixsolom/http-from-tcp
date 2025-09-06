package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	defer file.Close()

	for {
		buff := make([]byte, 8)
		n, err := file.Read(buff)
		if err == io.EOF {
			log.Printf("End of file reached: %v", err)
			break
		}
		if err != nil {
			log.Fatalf("unable to read from file: %v", err)
		}

		snippet := string(buff[:n])
		fmt.Printf("read: %s\n", snippet)
	}
	fmt.Println("Successfully finished reading the file")
}
