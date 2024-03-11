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

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		lines := strings.Split(text, "\n")

		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), "ping") {
				response += "+PONG\r\n"
			}
		}
	}

	buf := []byte(response)
	_, err := conn.Write(buf)
	if err != nil {
		return
	}
}
