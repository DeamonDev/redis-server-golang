package main

import (
	"bytes"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/domain/configuration"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type RedisServer struct {
	respParser         *resp.RespParser
	commandParser      *command.RedisCommandParser
	db                 map[string]DbRow
	mu                 *sync.RWMutex
	redisConfiguration configuration.RedisConfiguration
}

type DbRow struct {
	Value  string
	Expiry *time.Time
}

func NewRedisServer(configuration configuration.RedisConfiguration) *RedisServer {
	return &RedisServer{
		respParser:         resp.NewParser(),
		commandParser:      command.NewRedisCommandParser(configuration),
		db:                 make(map[string]DbRow),
		mu:                 &sync.RWMutex{},
		redisConfiguration: configuration,
	}
}

func main() {
	var port string
	port = "6379"

	var role string
	role = "master"

	var replicaOf *configuration.ReplicaOf
	replicaOf = nil

	for i, arg := range os.Args {
		switch arg {
		case "--port", "-p":
			port = os.Args[i+1]
		case "--replicaof":
			role = "slave"
			replicaOf = &configuration.ReplicaOf{
				Host: os.Args[i+1],
				Port: os.Args[i+2],
			}
		default:
			continue
		}

	}

	redisConfiguration := configuration.RedisConfiguration{
		Port:      port,
		Role:      role,
		ReplicaOf: replicaOf,
	}

	redisServer := NewRedisServer(redisConfiguration)

	l, err := net.Listen("tcp", "0.0.0.0"+":"+port)
	if err != nil {
		fmt.Printf("Failed to bind to the port %s", port)
		os.Exit(1)
	}

	fmt.Printf("Server is listening on port 6379 with configuration: %v", redisConfiguration)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error while accepting connection: ", err.Error())
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

		parsed, err := server.respParser.Parse(reader)
		if err != nil {
			fmt.Println(err)
			return
		}
		parsedCommand, err := server.commandParser.Parse(parsed)
		if err != nil {
			fmt.Println(err)
			return
		}

		switch commandValue := parsedCommand.(type) {
		case command.EchoCommand:
			str := fmt.Sprintf("+%s\r\n", commandValue.Value)
			conn.Write([]byte(str))
		case command.InfoCommand:
			rolePrefix := "role"
			role := commandValue.Role

			connectedSlavesPrefix := "connected_slaves"
			connectedSlaves := commandValue.ConnectedSlaves
			connectedSlavesNoOfDigits := len(strconv.Itoa(connectedSlaves))

			masterReplIdPrefix := "master_replid"
			masterReplId := commandValue.MasterReplId
			masterReplIdNoOfDigits := len(masterReplId)

			masterReplOffsetPrefix := "master_repl_offset"
			masterReplOffset := commandValue.MasterReplOffset
			masterReplOffsetNoOfDigits := len(strconv.Itoa(masterReplOffset))

			length := len(rolePrefix) + 1 + len(role) + 1 + len(connectedSlavesPrefix) + 1 + connectedSlavesNoOfDigits + 1 + len(masterReplIdPrefix) + 1 + masterReplIdNoOfDigits +
				1 + len(masterReplOffsetPrefix) + 1 + masterReplOffsetNoOfDigits

			str := fmt.Sprintf("$%d\r\n%s:%s\n%s:%d\n%s:%s\n%s:%d\r\n", length, rolePrefix, role, connectedSlavesPrefix, connectedSlaves, masterReplIdPrefix, masterReplId,
				masterReplOffsetPrefix, masterReplOffset)

			conn.Write([]byte(str))
		case command.SetCommand:
			server.mu.Lock()

			key := commandValue.Key
			value := commandValue.Value
			expiry := commandValue.Expiry

			server.db[key] = DbRow{Value: value, Expiry: expiry}
			server.mu.Unlock()

			str := "+OK\r\n"
			conn.Write([]byte(str))
		case command.GetCommand:
			var str string
			server.mu.Lock()

			key := commandValue.Key
			value, exists := server.db[key]
			length := len(value.Value)

			if exists {

				if value.Expiry != nil && time.Now().After(*value.Expiry) {
					str = "$-1\r\n"
				} else {
					str = fmt.Sprintf("$%d\r\n%s\r\n", length, value.Value)
				}
			} else {
				str = "$-1\r\n"
			}

			server.mu.Unlock()

			conn.Write([]byte(str))
		case command.PingCommand:
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}
