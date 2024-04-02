package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type CommandParserError struct {
	Domain   string
	Message  string
	Internal error
}

func (e *CommandParserError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Domain, e.Message, e.Internal)
}

func NewCommandParserError(message string, err error) *CommandParserError {
	return &CommandParserError{
		Domain:   "CommandParser",
		Message:  message,
		Internal: err,
	}
}

type RedisCommand interface{}

type EchoCommand struct {
	Value string
}

type PingCommand struct{}

type SetCommand struct {
	Key    string
	Value  string
	Expiry *time.Time
}

type GetCommand struct {
	Key string
}

type RedisCommandParser struct{}

func NewRedisCommandParser() *RedisCommandParser {
	return &RedisCommandParser{}
}

func (rcp *RedisCommandParser) Parse(respValue resp.RespValue) (RedisCommand, error) {

	switch s := respValue.(type) {
	case resp.ArrayRespValue:

		var args []resp.BulkStringRespValue
		for _, item := range s.Arr {
			if b, ok := item.(resp.BulkStringRespValue); ok {
				args = append(args, b)
			}
		}

		switch strings.ToLower(args[0].Str) {
		case "echo":
			return EchoCommand{Value: args[1].Str}, nil
		case "ping":
			return PingCommand{}, nil
		case "set":
			if len(args) > 3 {
				switch strings.ToLower(args[3].Str) {
				case "px":
					num, _ := strconv.Atoi(args[4].Str)
					expiry := time.Now().Add(time.Duration(num) * time.Millisecond)

					return SetCommand{Key: args[1].Str, Value: args[2].Str, Expiry: &expiry}, nil
				default:
					return nil, NewCommandParserError("unknown set command argument", nil)
				}
			} else {
				return SetCommand{Key: args[1].Str, Value: args[2].Str, Expiry: nil}, nil
			}
		case "get":
			return GetCommand{Key: args[1].Str}, nil
		default:
			return nil, NewCommandParserError("unknown command", nil)
		}

	default:
		return nil, NewCommandParserError("redis command should be represented as resp array of resp bulk strings", nil)
	}
}
