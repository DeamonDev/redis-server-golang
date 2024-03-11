package main

import (
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

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var response string

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	message := string(buffer[:n])
	fmt.Println(message)
	lines := strings.Split(message, "\n")

	for _, line := range lines {
		if strings.Contains(line, "ping") {
			response += "+PONG\r\n"
		}
	}

	buf := []byte(response)
	_, err = conn.Write(buf)
	if err != nil {
		return
	}
}
