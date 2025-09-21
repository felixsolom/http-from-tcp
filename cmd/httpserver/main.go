package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/felixsolom/http-from-tcp/internal/request"
	"github.com/felixsolom/http-from-tcp/internal/response"
	"github.com/felixsolom/http-from-tcp/internal/server"
)

const port = 42069

func newHandler(w io.Writer, req *request.Request) *server.HandlerError {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: response.BadRequest,
			Message:    "Your problem is not my problem\n",
		}
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: response.InternalServerError,
			Message:    "Woopsie, my bad\n",
		}
	}
	if err := response.WriteStatusLine(w, response.OK); err != nil {
		log.Printf("Failed to write status line: %v", err)
		return &server.HandlerError{
			StatusCode: response.InternalServerError,
			Message:    "Woopsie, my bad\n",
		}
	}
	body := "All good, frfr\n"
	if err := response.WriteHeaders(w, response.GetDefaultHeaders(len(body))); err != nil {
		log.Printf("Failed to write headers: %v", err)
		return &server.HandlerError{
			StatusCode: response.InternalServerError,
			Message:    "Woopsie, my bad\n",
		}
	}
	_, err := io.WriteString(w, body)
	if err != nil {
		log.Printf("Failed to write body: %v", err)
		return &server.HandlerError{
			StatusCode: response.InternalServerError,
			Message:    "Woopsie, my bad\n",
		}
	}
	return nil
}

func main() {
	server, err := server.Serve(port, newHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
