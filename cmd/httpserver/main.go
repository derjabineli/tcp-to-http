package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"net/http"
	"fmt"
	"errors"
	"io"

	"github.com/derjabineli/httpfromtcp/internal/server"
	"github.com/derjabineli/httpfromtcp/internal/request"
	"github.com/derjabineli/httpfromtcp/internal/response"

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
	w.WriteStatusLine(response.StatusOK)
	headers := response.GetDefaultHeaders(0)
	headers.Delete("content-length")
	headers.Set("Transfer-Encoding", "chunked")
	headers.Set("Host", "httpbin.org")
	w.WriteHeaders(headers)

	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
			w.WriteChunkedBodyDone()
			return
		}
		handler500(w)
		return
	}
	w.WriteChunkedBody(buf[:n])
	}
}
