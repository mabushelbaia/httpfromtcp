package main

import (
	"log/slog"
	"strconv"

	"github.com/mabushelbaia/httpfromtcp/internal/request"
	"github.com/mabushelbaia/httpfromtcp/internal/response"
	"github.com/mabushelbaia/httpfromtcp/internal/server"
)

func writeHTML(w *response.Writer, status response.StatusCode, html string) {
	// 1. write status line
	w.WriteStatusLine(status)

	hdrs := response.GetDefaultHeaders(0)
	hdrs.Replace("Content-Type", "text/html")
	hdrs.Replace("Content-Length", strconv.Itoa(len(html)))
	w.WriteHeaders(hdrs)
	w.WriteBody([]byte(html))
}

func main() {
	handler := func(w *response.Writer, req *request.Request) *server.HandlerError {
		slog.Info(req.RequestLine.RequestTarget)

		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			writeHTML(w, response.BadRequest, html400)
		case "/myproblem":
			writeHTML(w, response.ServerError, html500)
		default:
			writeHTML(w, response.Ok, html200)
		}

		// nil means no error, the response is already written
		return nil
	}

	_, err := server.Serve(42069, handler)
	if err != nil {
		panic(err)
	}

	select {} // block forever so server keeps running
}

// Example HTML constants:
const html400 = `<html>
  <head><title>400 Bad Request</title></head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

const html500 = `<html>
  <head><title>500 Internal Server Error</title></head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`

const html200 = `<html>
  <head><title>200 OK</title></head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
