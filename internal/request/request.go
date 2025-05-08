package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode"
)

type ParserState int

const (
  StateInitialized ParserState = iota
  StateDone
)

const bufferSize = 8

type Request struct {
  RequestLine RequestLine
  State ParserState
}

type RequestLine struct {
  Method        string
  RequestTarget string
  HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
  buffer := make([]byte, bufferSize)
  readToIndex := 0
  request := Request{
    State: StateInitialized,
  }

  for request.State != StateDone {
    if len(buffer) == cap(buffer) {
      newBuf := make([]byte, len(buffer) * 2)
      copy(newBuf, buffer)
      buffer = newBuf
    }

    read, readErr := reader.Read(buffer[readToIndex:])
    readToIndex += read

    if read > 0 {
      bytesParsed, parseErr := request.parse(buffer[:readToIndex])
      if parseErr != nil {
        if request.State == StateDone {
          break
        }
        return nil, parseErr
      }

      if bytesParsed > 0 {
        copy(buffer, buffer[bytesParsed:readToIndex])
        readToIndex -= bytesParsed
      }
    }

    if readErr == io.EOF {
      break
    } else if readErr != nil {
      return nil, readErr
    }
  }
  
  return &request, nil
}

func parseRequestLine(request *Request, data []byte) (int, error) {
  requestLineIndex := bytes.Index(data, []byte("\r\n"))
  if requestLineIndex == -1 {
    return 0, nil
  }

  requestLine := string(data[:requestLineIndex])
  parts := strings.Split(requestLine, " ")

  if len(parts) != 3 {
    return 0, errors.New("bad request line")
  }
  if !isUpper(parts[0]) {
    return 0, errors.New("invalid method")
  }
  if parts[2] != "HTTP/1.1" {
    return 0, errors.New("invalid http version")
  }

  request.RequestLine = RequestLine{
    HttpVersion: "1.1",
    Method: parts[0],
    RequestTarget: parts[1],
  }
  return requestLineIndex + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
  switch r.State {
  case StateInitialized:
    bytesParsed, err := parseRequestLine(r, data)
    if err != nil {
      return 0, err
    } 
    if bytesParsed == 0 {
      return 0, nil
    }
    r.State = StateDone
    return bytesParsed, nil
  case StateDone:
    return 0, errors.New("error: trying to read data in a done state")
  default:
    return 0, errors.New("error: unknown state")
  }
}


func isUpper(s string) bool {
  for _, r := range s {
    if !unicode.IsUpper(r) && unicode.IsLetter(r) {
      return false
    }
  }
  return true
}

