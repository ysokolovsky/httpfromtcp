package main

import (
	"fmt"
	"log"
	"net"

	"httpfromtcp/internal/request"
)

func main() {

	listener, _ := net.Listen("tcp", ":42069")
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		fmt.Println("conn accepted", conn)

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("%s", err)
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range r.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Println(string(r.Body))

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}

}
