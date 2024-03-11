package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	fmt.Println("Server is listening on port 6379")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var response string

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // End of input
		}

		if strings.Contains(line, "ping") {
			response += "PONG\n"
		}
	}

	buf := []byte(response)
	_, err := conn.Write(buf)
	if err != nil {
		return
	}

}
