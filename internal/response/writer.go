package response

import (
	"errors"
	"fmt"

	"github.com/derjabineli/httpfromtcp/internal/headers"
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
  if w.State!= WriteStatusLine {
		return errors.New("writing status line out of order")
	}
	
	statusLine := getStatusLine(statusCode)
	_, err := w.Writer.Write(statusLine)
	w.State = WriteHeaders
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
  if w.State!= WriteHeaders {
		return errors.New("writing headers out of order")
	}
	writableHeader := []byte{}
  for header, value := range headers {
    writableHeader = fmt.Appendf(writableHeader, fmt.Sprintf("%v: %v\r\n", header, value)) 
  }
  writableHeader = fmt.Appendf(writableHeader, "\r\n")
  _, err := w.Writer.Write(writableHeader)
	w.State = WriteBody
  return err
}

func (w *Writer) WriteBody(b []byte) (int, error) {
	if w.State != WriteBody {
		return 0, errors.New("writing body out of order")
	}
	n, err := w.Writer.Write(b)
	return n, err
}