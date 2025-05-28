package response

import (
  "io"
  "errors"
  "strconv"
  "fmt"

  "github.com/derjabineli/httpfromtcp/internal/headers"
)

type StatusCode int

const (
  StatusOK StatusCode = 200
  StatusBadRequest = 400
  StatusInternalServerError = 500
)

type WriterState int

const (
	WriteStatusLine WriterState = iota
	WriteHeaders 
	WriteBody 
)

type Writer struct {
	state WriterState
	Writer io.Writer	
}

func getStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""
	switch statusCode {
	case StatusOK:
		reasonPhrase = "OK"
	case StatusBadRequest:
		reasonPhrase = "Bad Request"
	case StatusInternalServerError:
		reasonPhrase = "Internal Server Error"
	}
	statusLine := []byte(fmt.Sprintf("HTTP/1.1 %v %v \r\n", statusCode, reasonPhrase))
	return statusLine
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
  if w.state != WriteStatusLine {
		return errors.New("writing status line out of order")
	}
	
	statusLine := getStatusLine(statusCode)
	_, err := w.Writer.Write(statusLine)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
  h := headers.NewHeaders()
  h.SetHeader("Content-Length", strconv.Itoa(contentLen))
  h.SetHeader("Connection", "close")
  h.SetHeader("Content-Type", "text/plain")
  return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
  if w.state != WriteHeaders {
		return errors.New("writing headers out of order")
	}
	writableHeader := []byte{}
  for header, value := range headers {
    writableHeader = fmt.Appendf(writableHeader, fmt.Sprintf("%v: %v\r\n", header, value)) 
  }
  writableHeader = fmt.Appendf(writableHeader, "\r\n")
  _, err := w.Writer.Write(writableHeader)
  return err
}

func (w *Writer) WriteBody(b []byte) error {
	if w.state != WriteBody {
		return errors.New("writing body out of order")
	}
	return nil
}
