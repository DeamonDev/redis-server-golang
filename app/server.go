package main

import (
	"bytes"
	"fmt"
	command "github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/parser"
	"net"
	"os"
)

type RedisServer struct {
	respParser    *parser.RespParser
	commandParser *command.RedisCommandParser
}

func NewRedisServer() *RedisServer {
	return &RedisServer{
		respParser:    parser.NewParser(),
		commandParser: command.NewRedisCommandParser(),
	}
}

func main() {
	redisServer := NewRedisServer()

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

		go handleClient(conn, redisServer)
	}
}

func handleClient(conn net.Conn, server *RedisServer) {
	defer conn.Close()

	buff := make([]byte, 1024)
	for {
		_, err := conn.Read(buff) // we assume messages are short. As for now
		if err != nil {
			return
		}

		reader := bytes.NewReader(buff)

		parsed, _ := server.respParser.Parse(reader)
		parsedCommand, _ := server.commandParser.Parse(parsed)

		switch commandValue := parsedCommand.(type) {
		case command.EchoCommand:
			str := fmt.Sprintf("+%s\r\n", commandValue.Value)
			conn.Write([]byte(str))
		case command.PingCommand:
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}
