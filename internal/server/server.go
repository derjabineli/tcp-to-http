package server

import (
	"fmt"
	"net"
  "sync/atomic"
  "log"
	"io"
	"bytes"

  "github.com/derjabineli/httpfromtcp/internal/response"
  "github.com/derjabineli/httpfromtcp/internal/request"

)

type Server struct {
  listener net.Listener 
  closed atomic.Bool
	handler Handler
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode 
	Message string
}

func (h HandlerError)writeError(w io.Writer) {
	response.WriteStatusLine(w, h.StatusCode)
	headers := response.GetDefaultHeaders(len(h.Message))
	response.WriteHeaders(w, headers)
	w.Write([]byte(h.Message))
}


func Serve(port int, handler Handler) (*Server, error) {
  portAddr := fmt.Sprintf(":%d", port)
  listener, err := net.Listen("tcp", portAddr)
  if err != nil {
    return nil, err
  }

  s := &Server{
    listener: listener,
		handler: handler,
  }

  go s.listen()
  return s, nil
}

func (s *Server) Close() error {
  s.closed.Store(true)
  if s.listener != nil {
    return s.listener.Close()
  }
  return nil
}

func (s *Server) listen() {
  for {
    conn, err := s.listener.Accept()
    if err != nil {
      if s.closed.Load() {
        return
      }
      log.Printf("Couldn't accept connection: %v", err)
      continue
    }

    go s.handle(conn) 
  }
 }

func (s *Server) handle(conn net.Conn) {
  defer conn.Close() 
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message: err.Error(),
		}
		hErr.writeError(conn)
		return
	}
	buf := bytes.NewBuffer([]byte{}) 
	hErr := s.handler(buf, req)
	if hErr != nil {
		hErr.writeError(conn)
		return
	}
  response.WriteStatusLine(conn, response.StatusOK)
  headers := response.GetDefaultHeaders(buf.Len())
  response.WriteHeaders(conn, headers)
	buf.WriteTo(conn)
}
