package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/derjabineli/httpfromtcp/internal/request"
	"github.com/derjabineli/httpfromtcp/internal/response"
)

type Server struct {
  listener net.Listener 
  closed atomic.Bool
	handler Handler
}

type Handler func(w *response.Writer, req *request.Request)


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
	w := &response.Writer{
		State: response.WriteStatusLine,
		Writer: conn,
	}
	req, err := request.RequestFromReader(conn)
  if err != nil {
    w.WriteStatusLine(response.StatusBadRequest)
    body := []byte(fmt.Sprintf("Error parsing request %v", err))
    w.WriteHeaders(response.GetDefaultHeaders(len(body)))
    w.WriteBody(body)
  }
	s.handler(w, req)
}
