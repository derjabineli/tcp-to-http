package response

import (
	"strconv"

	"github.com/derjabineli/httpfromtcp/internal/headers"
)

type StatusCode int

const (
  StatusOK 					StatusCode = 200
  StatusBadRequest 			StatusCode = 400
  StatusInternalServerError StatusCode = 500
)

func GetDefaultHeaders(contentLen int) headers.Headers {
  h := headers.NewHeaders()
  h.SetHeader("Content-Length", strconv.Itoa(contentLen))
  h.SetHeader("Connection", "close")
  h.SetHeader("Content-Type", "text/plain")
  return h
}
