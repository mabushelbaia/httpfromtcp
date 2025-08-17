package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mabushelbaia/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen(
		"tcp",
		":42069",
	)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	for {

		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())
		r, err := request.RequestFromReader(conn)

		if err != nil {
			log.Fatal("error", "error", err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
		fmt.Println("Headers:")
		r.Headers.ForEach(func(key, value string) {
			fmt.Printf("- %s: %s\n", key, value)
		})

		fmt.Println("Body:")
		fmt.Println(string(r.Body))

	}
}
