package headers

import (
	"bytes"
	"errors"
	"strings"
  "unicode"
)

const crlf = "\r\n"

// Ensure Range Table is properly sorted
var validHttpTokenRunes = &unicode.RangeTable{ 
	R16: []unicode.Range16{
	  {Lo: '!', Hi: '!', Stride: 1},
		{Lo: '#', Hi: '#', Stride: 1},
		{Lo: '$', Hi: '$', Stride: 1},
		{Lo: '%', Hi: '%', Stride: 1},
		{Lo: '&', Hi: '&', Stride: 1},
		{Lo: '\'', Hi: '\'', Stride: 1},
		{Lo: '*', Hi: '*', Stride: 1},
		{Lo: '+', Hi: '+', Stride: 1},
		{Lo: '-', Hi: '-', Stride: 1}, 
		{Lo: '.', Hi: '.', Stride: 1},
		{Lo: '0', Hi: '9', Stride: 1},
		{Lo: 'A', Hi: 'Z', Stride: 1},
		{Lo: '^', Hi: '^', Stride: 1},
		{Lo: '_', Hi: '_', Stride: 1},
		{Lo: '`', Hi: '`', Stride: 1},
    {Lo: 'a', Hi: 'z', Stride: 1},
		{Lo: '|', Hi: '|', Stride: 1},
		{Lo: '~', Hi: '~', Stride: 1},	
  },
}

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

  if !isValidTChar(fieldName) {
    return 0, false, errors.New("contains invalid runes")
  } 

  h.SetHeader(fieldName, fieldValue)

  return idx + 2, false, nil
}

func (h Headers) SetHeader(fieldName, fieldValue string) {
  loweredName := strings.ToLower(fieldName)
  name, exists := h[loweredName]
  if exists {
    h[loweredName] = name + ", " + fieldValue
    return
  }
  h[loweredName] = fieldValue
}

func isValidTChar(tChar string) bool {
  for _, c := range tChar {
    if !unicode.Is(validHttpTokenRunes, rune(c)) {
      return false
    }
  }
  return true
}
