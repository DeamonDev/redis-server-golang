package command

import (
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

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

func (rcp *RedisCommandParser) Parse(respValue parser.RespValue) (RedisCommand, error) {

	switch s := respValue.(type) {
	case parser.ArrayRespValue:

		var args []parser.BulkStringRespValue
		for _, item := range s.Arr {
			if b, ok := item.(parser.BulkStringRespValue); ok {
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
					return nil, nil
				}
			} else {
				return SetCommand{Key: args[1].Str, Value: args[2].Str, Expiry: nil}, nil
			}
		case "get":
			return GetCommand{Key: args[1].Str}, nil
		default:
			return nil, nil
		}

	default:
		return nil, nil
	}
}
