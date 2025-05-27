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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
  switch statusCode {
  case StatusOK:
    _, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
    return err
  case StatusBadRequest:
    _, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
    return err
  case StatusInternalServerError:
    _, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
    return err
  default:
    return errors.New("bad status code")
  }
}

func GetDefaultHeaders(contentLen int) headers.Headers {
  h := headers.NewHeaders()
  h.SetHeader("Content-Length", strconv.Itoa(contentLen))
  h.SetHeader("Connection", "close")
  h.SetHeader("Content-Type", "text/plain")
  return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
  writableHeader := []byte{}
  for header, value := range headers {
    writableHeader = fmt.Appendf(writableHeader, fmt.Sprintf("%v: %v\r\n", header, value)) 
  }
  writableHeader = fmt.Appendf(writableHeader, "\r\n")
  _, err := w.Write(writableHeader)
  return err
}
