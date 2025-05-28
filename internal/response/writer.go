package response

import (
	"errors"
	"fmt"
	"io"

	"github.com/derjabineli/httpfromtcp/internal/headers"
)

type WriterState int

const (
	writerStateStatusLine WriterState = iota
	writerStateHeaders 
	writerStateBody 
)

type Writer struct {
	state WriterState
	Writer io.Writer	
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		state: writerStateStatusLine,
		Writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
  if w.state!= writerStateStatusLine {
		return errors.New("writing status line out of order")
	}
	
	statusLine := getStatusLine(statusCode)
	_, err := w.Writer.Write(statusLine)
	w.state = writerStateHeaders
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
  	if w.state!= writerStateHeaders {
		return errors.New("writing headers out of order")
	}
  	for header, value := range headers {
		w.Writer.Write([]byte(fmt.Sprintf("%v: %v\r\n", header, value)))
  	}
  	_, err := w.Writer.Write([]byte("\r\n"))
	w.state = writerStateBody
  return err
}

func (w *Writer) WriteBody(b []byte) (int, error) {
	if w.state != writerStateBody {
		return 0, errors.New("writing body out of order")
	}
	n, err := w.Writer.Write(b)
	return n, err
}