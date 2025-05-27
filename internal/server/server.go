package server

import (
	"fmt"
	"net"
  "sync/atomic"
  "log"

  "github.com/derjabineli/httpfromtcp/internal/response"
)

type Server struct {
  listener net.Listener 
  closed atomic.Bool
}

func Serve(port int) (*Server, error) {
  portAddr := fmt.Sprintf(":%d", port)
  listener, err := net.Listen("tcp", portAddr)
  if err != nil {
    return nil, err
  }

  s := &Server{
    listener: listener,
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

    go s.handler(conn) 
  }
 }

func (s *Server) handler(conn net.Conn) {
  defer conn.Close()  
  err := response.WriteStatusLine(conn, 200)
  if err != nil {
    return
  }
  headers := response.GetDefaultHeaders(0)
  err = response.WriteHeaders(conn, headers)
  if err != nil {
    return
  }
}
