package response

import (
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/mabushelbaia/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Ok          StatusCode = 200
	BadRequest  StatusCode = 400
	ServerError StatusCode = 500
)

type WriterState int

const (
	StatusLine WriterState = iota
	Headers
	Body
)

type Writer struct {
	State WriterState
	Conn  net.Conn
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.State != StatusLine {
		return fmt.Errorf("status line already written or wrong order")
	}
	fmt.Fprintf(w.Conn, "HTTP/1.1 %d %s\r\n", statusCode, statusText(statusCode))
	w.State++
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.State != Headers {
		return fmt.Errorf("headers line already written or wrong order")
	}
	headers.ForEach(func(key, value string) {
		fmt.Fprintf(w.Conn, "%s: %s\r\n", key, value)
	})
	fmt.Fprint(w.Conn, "\r\n")
	w.State++
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.State != Body {
		return 0, fmt.Errorf("body line already written or wrong order")
	}
	_, err := w.Conn.Write(p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func statusText(statusCode StatusCode) string {
	switch statusCode {
	case Ok:
		return "OK"
	case BadRequest:
		return "Bad Request"
	case ServerError:
		return "Internal Server Error"
	default:
		return "Internal Server Error"
	}
}
func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var line string
	switch statusCode {
	case Ok:
		line = "HTTP/1.1 200 OK\r\n"
	case BadRequest:
		line = "HTTP/1.1 400 Bad Request\r\n"
	case ServerError:
		line = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		line = ""
	}
	_, err := w.Write([]byte(line))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("content-length", strconv.Itoa(contentLen))
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	h.ForEach(func(key, value string) {
		line := key + ": " + value + "\r\n"
		w.Write([]byte(line))
	})

	// Write the blank line after headers
	_, err := w.Write([]byte("\r\n"))
	return err
}
