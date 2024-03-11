package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	fmt.Println("Server is listening on port 6379")

	for {
		_, err = l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	buf := []byte("+PONG\r\n")
	_, err := conn.Write(buf)
	if err != nil {
		return
	}
}
