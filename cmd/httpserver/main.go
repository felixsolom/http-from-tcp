package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/felixsolom/http-from-tcp/internal/request"
	"github.com/felixsolom/http-from-tcp/internal/response"
	"github.com/felixsolom/http-from-tcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
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

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	handler200(w, req)
	return
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.BadRequest)
	body := []byte(
		`
			<html>
				<head>
					<title>400 Bad Request</title>
				</head>
				<body>
					<h1>Bad Request</h1>
					<p>Your request honestly kinda sucked.</p>
				</body>
				</html>`,
	)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.InternalServerError)
	body := []byte(
		`
		<html>
			<head>
				<title>500 Internal Server Error</title>
			</head>
			<body>
				<h1>Internal Server Error</h1>
				<p>Okay, you know what? This one is on me.</p>
			</body>
			</html>`,
	)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.OK)
	body := []byte(
		`
		<html>
			<head>
				<title>200 OK</title>
			</head>
			<body>
				<h1>Success!</h1>
				<p>Your request was an absolute banger.</p>
			</body>
			</html>`,
	)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func proxyHandler(w *response.Writer, req *request.Request) {
	proxyPrefex := "/httpbin/"
	targetServer := "https://httpbin.org/"

	if !strings.HasPrefix(req.RequestLine.RequestTarget, proxyPrefex) {
		handler400(w, req)
		return
	}
	targetURL := targetServer + strings.TrimPrefix(req.RequestLine.RequestTarget, proxyPrefex)
	proxyRes, err := http.Get(targetURL)
	if err != nil {
		log.Printf("Couldn't get a response from http_bin: %v", err)
		return
	}
	defer proxyRes.Body.Close()

	w.WriteStatusLine(response.StatusCode(proxyRes.StatusCode))
	h := response.GetDefaultHeaders(0)
	for key, value := range proxyRes.Header {
		h.Set(key, strings.Join(value, ","))
	}
	delete(h, "Content-Length")
	h.Set("Transfer-Encoding", "chunked")
	w.WriteHeaders(h)

	for {
		buf := make([]byte, 1024)
		n, err := proxyRes.Body.Read(buf)
		if n > 0 {
			_, err = w.WriteChunkedBody(buf[:n])
			if err != nil {
				log.Printf("Couldn't write chunk to body: %v", err)
				break
			}
			go func() { log.Printf("read: %d bytes\n", n) }()
		}
		if err == io.EOF {
			log.Printf("reached EOF")
			break
		}

		if err != nil {
			log.Printf("Coudn't read response body: %v", err)
			break
		}

	}
	if _, err = w.WriteChunkedBodyDone(); err != nil {
		log.Printf("Couldn't write end chunk to response: %v", err)
	}
	log.Printf("finished writing to response body")
}
