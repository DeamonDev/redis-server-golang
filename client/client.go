package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	message := "*2\r\n$4\r\nECHO\r\n$3\r\nstrawberry\r\n"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing:", err.Error())
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	// Print received response
	fmt.Printf("Received response: %s\n", buffer[:n])
}
