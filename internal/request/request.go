package request

import (
  "io"
  "fmt"
  "strings"
  "errors"
  "unicode"
)

type Request struct {
  RequestLine RequestLine
}

type RequestLine struct {
  HttpVersion   string
  RequestTarget string
  Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
  req, err := io.ReadAll(reader)
  if err != nil {
    fmt.Printf("Couldn't read request. Err %v\n", err.Error())
    return nil, err
  }
  
  requestLine, err := parseRequestLine(req)
  if err != nil {
    return nil, err
  }

  return &Request{
    RequestLine: requestLine,
  }, nil
}

func parseRequestLine(req []byte) (RequestLine, error) {
  requestParts := strings.Split(string(req), "\r\n")
  parts := strings.Split(requestParts[0], " ")
  if len(parts) != 3 {
    return RequestLine{}, errors.New("bad request line")
  }
  
  if !isUpper(parts[0]) {
    return RequestLine{}, errors.New("invalid method")
  }
  if parts[2] != "HTTP/1.1" {
    return RequestLine{}, errors.New("invalid http version")
  }

  return RequestLine{
    HttpVersion: "1.1",
    RequestTarget: parts[1],
    Method: parts[0],
  }, nil
}


func isUpper(s string) bool {
  for _, r := range s {
    if !unicode.IsUpper(r) && unicode.IsLetter(r) {
      return false
    }
  }
  return true
}
