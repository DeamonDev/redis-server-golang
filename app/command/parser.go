package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/parser"
	"strings"
)

type RedisCommand interface{}

type EchoCommand struct {
	Value string
}

type PingCommand struct{}

type SetCommand struct {
	Key   string
	Value string
}

type GetCommand struct {
	Key string
}

type RedisCommandParser struct{}

func NewRedisCommandParser() *RedisCommandParser {
	return &RedisCommandParser{}
}

func (rcp *RedisCommandParser) Parse(respValue parser.RespValue) (RedisCommand, error) {

	switch s := respValue.(type) {
	case parser.ArrayRespValue:

		var args []parser.BulkStringRespValue
		for _, item := range s.Arr {
			if b, ok := item.(parser.BulkStringRespValue); ok {
				args = append(args, b)
			}
		}

		switch strings.ToLower(args[0].BulkStr) {
		case "echo":
			return EchoCommand{Value: args[1].BulkStr}, nil
		case "ping":
			return PingCommand{}, nil
		case "set":
			return SetCommand{Key: args[1].BulkStr, Value: args[2].BulkStr}, nil
		case "get":
			return GetCommand{Key: args[1].BulkStr}, nil
		default:
			return nil, nil
		}

	default:
		return nil, nil
	}
}
