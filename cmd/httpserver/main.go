package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"net/http"
	"fmt"
	"io"
	"crypto/sha256"

	"github.com/derjabineli/httpfromtcp/internal/server"
	"github.com/derjabineli/httpfromtcp/internal/request"
	"github.com/derjabineli/httpfromtcp/internal/response"

)

const port = 42010

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
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w)
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		httpBinProxy(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/video" {
		handlerVideo(w, req)
		return
	}
	handler200(w)
	return
}

func handler400(w *response.Writer) {
	body := []byte(`
		<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
	`)
	w.WriteStatusLine(response.StatusBadRequest)
	headers:= response.GetDefaultHeaders(len(body))
	headers.Overwrite("content-type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(body)
	return
}

func handler500(w *response.Writer) {
	body := []byte(`
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
	`)
	w.WriteStatusLine(response.StatusInternalServerError)
	headers:= response.GetDefaultHeaders(len(body))
	headers.Overwrite("content-type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(body)
	return
}

func handler200(w *response.Writer) {
	body := []byte(`
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
	`)
	w.WriteStatusLine(response.StatusOK)
	headers:= response.GetDefaultHeaders(len(body))
	headers.Overwrite("content-type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(body)
	return
}

func httpBinProxy(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := fmt.Sprintf("https://httpbin.org/%s", target)
	resp, err := http.Get(url)
	if err != nil {
		handler500(w)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusOK)
	headers := response.GetDefaultHeaders(0)
	headers.Delete("content-length")
	headers.Overwrite("Transfer-Encoding", "chunked")
	headers.Overwrite("Trailers", "X-Content-SHA256, X-Content-Length")
	w.WriteHeaders(headers)

	fullBody := []byte{}
	var maxBufferSize = 1024
	buf := make([]byte, maxBufferSize)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, err = w.WriteChunkedBody(buf[:n])
			if err != nil {
				fmt.Println("Error writing chunked body: ", err)
				break
			}
			fullBody = append(fullBody, buf[:n]...)
		}
		if err == io.EOF {
			break	
		}
		if err != nil {
			fmt.Println("Error reading from response body: ", err)
			break
		}
	}
	w.WriteChunkedBodyDone()

	sha256 := fmt.Sprintf("%x", sha256.Sum256(fullBody))
	trailers := response.GetDefaultHeaders(0)
	trailers.Delete("content-length")
	trailers.Overwrite("X-Content-SHA256", sha256)
	trailers.Overwrite("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
	w.WriteTrailers(trailers)
}

func handlerVideo(w *response.Writer, req *request.Request) {
	body, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		handler500(w)
		return
	}

	w.WriteStatusLine(response.StatusOK)
	headers := response.GetDefaultHeaders(len(body))
	headers.Overwrite("Content-Type", "video/mp4")
	w.WriteHeaders(headers)
	w.WriteBody(body)
}
