package response

import (
	"io"
	"strconv"

	"github.com/mabushelbaia/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Ok          StatusCode = 200
	BadRequest  StatusCode = 400
	ServerError StatusCode = 500
)

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
