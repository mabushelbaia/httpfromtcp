package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/mabushelbaia/httpfromtcp/internal/headers"
)

var ErrorBadRequestLine = fmt.Errorf("bad request line")
var ErrorBadMethod = fmt.Errorf("bad method")
var sep []byte = []byte("\r\n")

type RequestState int

const (
	InitState RequestState = iota
	HeadersState
	BodyState
	DoneState
	ErrorStaate
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       RequestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}
type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}
func (r *Request) Parse(b []byte) (int, error) {
	read := 0

outer:
	for {
		switch r.state {
		case InitState:
			rl, n, err := parseRequestLine(b[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer // need more data
			}
			r.RequestLine = *rl
			read += n
			r.state++

		case HeadersState:
			n, done, err := r.Headers.Parse(b[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			read += n
			if done {
				r.state++
			}

		case BodyState:
			val, ok := r.Headers.Get("Content-Length")
			if !ok {
				// No Content-Length, slurp everything and finish
				r.Body = append(r.Body, b[read:]...)
				read = len(b)
				r.state++
				break outer
			}

			bodySize, err := strconv.Atoi(val)
			if err != nil {
				return 0, fmt.Errorf("invalid Content-Length: %w", err)
			}

			remaining := min(bodySize-len(r.Body), len(b[read:]))
			if remaining <= 0 {
				return 0, fmt.Errorf("body overflow")
			}

			r.Body = append(r.Body, b[read:read+remaining]...)
			read += remaining

			if len(r.Body) == bodySize {
				r.state++
			} else {
				break outer // need more data
			}

		case DoneState:
			break outer
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.state == DoneState || r.state == ErrorStaate
}
func newRequest() *Request {
	return &Request{
		state:   InitState,
		Headers: headers.NewHeaders(),
	}
}
func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buff := make([]byte, 1)

	bufLen := 0

	for !request.done() {

		if bufLen >= len(buff) {
			newBuff := make([]byte, len(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}

		n, err := reader.Read(buff[bufLen:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.state = DoneState
				break
			}
			return nil, err
		}
		bufLen += n
		readN, err := request.Parse(buff[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buff, buff[readN:bufLen])
		bufLen -= readN
	}
	return request, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, sep)

	if idx == -1 {
		return nil, 0, nil
	}

	requestLine := b[:idx]
	read := idx + len(sep)

	parts := bytes.Split(requestLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrorBadRequestLine
	}
	method, target, version := parts[0], parts[1], parts[2]

	version = bytes.Split(version, []byte("/"))[1]

	if !bytes.Equal(bytes.ToUpper(method), method) {
		return nil, 0, ErrorBadMethod
	}

	return &RequestLine{
		Method:        string(method),
		RequestTarget: string(target),
		HttpVersion:   string(version),
	}, read, nil
}
