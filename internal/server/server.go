package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/mabushelbaia/httpfromtcp/internal/request"
	"github.com/mabushelbaia/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request) *HandlerError

func Serve(port uint16, h Handler) (*Server, error) {

	ln, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%d", port),
	)
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: ln,
		handler:  h,
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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.BadRequest,
			Message:    err.Error(),
		}
		hErr.write(conn)
		return
	}
	writer := response.Writer{
		Conn:  conn,
		State: response.StatusLine,
	}

	hErr := s.handler(&writer, req)

	if hErr != nil {
		hErr.write(conn)
		return
	}

}

func (hErr *HandlerError) write(w io.Writer) {
	// status line
	response.WriteStatusLine(w, hErr.StatusCode)

	// headers
	body := []byte(hErr.Message)
	hdrs := response.GetDefaultHeaders(len(body))
	response.WriteHeaders(w, hdrs)

	// body
	w.Write(body)
}
