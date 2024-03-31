package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	command "github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

func ExpiryAnalyzer(db map[string]DbRow, mu *sync.RWMutex) {
	for {
		fmt.Println("<<EXPIRY_ANALYZER>>")

		for k, v := range db {
			if v.Expiry == nil {
				continue
			}

			if time.Now().After(*v.Expiry) {
				mu.Lock()
				delete(db, k)
				mu.Unlock()
			}
		}

		time.Sleep(200 * time.Millisecond)
	}
}

type RedisServer struct {
	respParser    *parser.RespParser
	commandParser *command.RedisCommandParser
	db            map[string]DbRow
	mu            *sync.RWMutex
}

type DbRow struct {
	Value  string
	Expiry *time.Time
}

func NewRedisServer() *RedisServer {
	return &RedisServer{
		respParser:    parser.NewParser(),
		commandParser: command.NewRedisCommandParser(),
		db:            make(map[string]DbRow),
		mu:            &sync.RWMutex{},
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

	go ExpiryAnalyzer(redisServer.db, redisServer.mu)

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
		case command.SetCommand:
			str := "+OK\r\n"
			server.mu.Lock()

			key := commandValue.Key
			value := commandValue.Value
			expiry := commandValue.Expiry

			server.db[key] = DbRow{Value: value, Expiry: expiry}
			server.mu.Unlock()

			conn.Write([]byte(str))
		case command.GetCommand:
			server.mu.Lock()

			key := commandValue.Key
			value := server.db[key]
			length := len(value.Value)

			str := fmt.Sprintf("$%d\r\n%s\r\n", length, value)
			server.mu.Unlock()

			conn.Write([]byte(str))
		case command.PingCommand:
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}
