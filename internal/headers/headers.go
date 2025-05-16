package headers

import (
	"bytes"
	"errors"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
  idx := bytes.Index(data, []byte(crlf))
  if idx == -1 {
    return 0, false, nil
  }
  if idx == 0 {
    return 2, true, nil
  }

  parts := bytes.SplitN(data[:idx], []byte(":"), 2)
  fieldName := string(parts[0])

  if fieldName != strings.TrimRight(fieldName, " ") {
    return 0, false, errors.New("invalid header field name")
  }

  fieldName = strings.TrimSpace(fieldName)
  fieldValue := strings.TrimSpace(string(parts[1]))
  h.SetHeader(fieldName, fieldValue)

  return idx + 2, false, nil
}

func (h Headers) SetHeader(fieldName, fieldValue string) {
  h[fieldName] = fieldValue
}
