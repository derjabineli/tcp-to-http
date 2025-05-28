package response

import (
	"io"
	"strconv"

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
	State WriterState
	Writer io.Writer	
}

func GetDefaultHeaders(contentLen int) headers.Headers {
  h := headers.NewHeaders()
  h.SetHeader("Content-Length", strconv.Itoa(contentLen))
  h.SetHeader("Connection", "close")
  h.SetHeader("Content-Type", "text/plain")
  return h
}
