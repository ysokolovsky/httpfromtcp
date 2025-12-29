package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	udp, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
	fmt.Println("error ResolveUDPAddr", err)
	}

	conn, err := net.DialUDP("udp", nil, udp)
	if err != nil {
	fmt.Println("error DialUDP", err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error ReadString", err)
		}

		_, err = conn.Write([]byte(input))
		if err != nil {
			fmt.Println("error sending udp", err)
		}

	}

}
