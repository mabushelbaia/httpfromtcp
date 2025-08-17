package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/mabushelbaia/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {

	ln, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%d", port),
	)
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: ln,
	}
	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return // exit gracefully
			}
			log.Printf("accept error: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	body := []byte("{Hello World!}\r\n")
	hdrs := response.GetDefaultHeaders(len(body))
	response.WriteStatusLine(conn, response.Ok)
	response.WriteHeaders(conn, hdrs)
	conn.Write(body)
}
