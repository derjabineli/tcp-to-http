package request

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/derjabineli/httpfromtcp/internal/headers"
)

type ParserState int

const (
  requestStateInitialized ParserState = iota
  requestStateParsingHeaders
  requestStateParsingBody
  requestStateDone
)

const bufferSize = 8

type Request struct {
  RequestLine RequestLine
  Headers headers.Headers
  State ParserState
  Body []byte
}

type RequestLine struct {
  Method        string
  RequestTarget string
  HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
  buffer := make([]byte, bufferSize)
  readToIndex := 0
  request := &Request{
    State: requestStateInitialized,
    Headers: headers.NewHeaders(),
  }
  for request.State != requestStateDone {
    if readToIndex >= len(buffer) {
      newBuf := make([]byte, len(buffer) * 2)
      copy(newBuf, buffer)
      buffer = newBuf
    }

    numBytesRead, err := reader.Read(buffer[readToIndex:])
    if err != nil {
      if errors.Is(err, io.EOF) {
        if request.State != requestStateDone {
          return nil, errors.New("incomplete request")
        }
        break
      }
    }

    readToIndex += numBytesRead
    numBytesParsed, err := request.parse(buffer[:readToIndex])
    if err != nil {
      return nil, err
    }

    if numBytesParsed > 0 {
      copy(buffer, buffer[numBytesParsed:])
      readToIndex -= numBytesParsed
    }
  }
  
  return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
  totalBytesParsed := 0
  for r.State != requestStateDone {
    n, err := r.parseSingle(data[totalBytesParsed:])
    if err != nil {
      return 0, err
    }
    totalBytesParsed += n
    if n == 0 {
      break
    }
  }
  return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
  switch r.State {
  case requestStateInitialized:
    bytesParsed, err := parseRequestLine(r, data)
    if err != nil {
      return 0, err
    } 
    if bytesParsed == 0 {
      return 0, nil
    }
    r.State = requestStateParsingHeaders 
    return bytesParsed, nil
  case requestStateParsingHeaders:
    n, done, err := r.Headers.Parse(data)
    if err != nil {
      return 0, err
    }
    if done {
      r.State = requestStateParsingBody
    }
    return n, nil
  case requestStateParsingBody:
    headerValue, err := r.Headers.Get("Content-Length")
    if err != nil {
      r.State = requestStateDone
      return 0, nil
    }
    contentLength, err := strconv.Atoi(headerValue)
    if err != nil {
      return 0, errors.New("malformed content-length header")
    }

    r.Body = append(r.Body, data...)
    if len(r.Body) > contentLength {
      return 0, errors.New("request body size exceeds content length")
    } 
    if len(r.Body) == contentLength {
      r.State = requestStateDone
    }
    return len(data), nil
  case requestStateDone:
    return 0, errors.New("error: trying to read data in a done state")
  default:
    return 0, errors.New("error: unknown state")
  }
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

func isUpper(s string) bool {
  for _, r := range s {
    if !unicode.IsUpper(r) && unicode.IsLetter(r) {
      return false
    }
  }
  return true
}
