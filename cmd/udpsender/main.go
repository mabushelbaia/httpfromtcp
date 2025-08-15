package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	listner, err := net.ResolveUDPAddr(
		"udp",
		"localhost:42069",
	)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, listner)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		line, _ := reader.ReadString('\n')
		conn.Write([]byte(line))

	}

}
