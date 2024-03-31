package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"os"
	"testing"
	"time"
)

var cp *RedisCommandParser

func TestMain(m *testing.M) {
	cp = NewRedisCommandParser()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestPingCommand(t *testing.T) {
	pingResp := resp.ArrayRespValue{
		Arr: []resp.RespValue{
			resp.BulkStringRespValue{
				Str: "PING",
			},
		},
	}

	pingCommand, _ := cp.Parse(pingResp)
	expected := PingCommand{}

	if pingCommand != expected {
		t.Errorf("should be PING command")
	}
}

func TestEchoCommand(t *testing.T) {
	echoResp := resp.ArrayRespValue{
		Arr: []resp.RespValue{
			resp.BulkStringRespValue{
				Str: "ECHO",
			},
			resp.BulkStringRespValue{
				Str: "FOO",
			},
		},
	}

	echoCommand, _ := cp.Parse(echoResp)
	expected := EchoCommand{
		Value: "FOO",
	}

	if echoCommand != expected {
		t.Errorf("Should be echo command")
	}
}

func TestSetCommand(t *testing.T) {
	setResp := resp.ArrayRespValue{
		Arr: []resp.RespValue{
			resp.BulkStringRespValue{
				Str: "SET",
			},
			resp.BulkStringRespValue{
				Str: "orange",
			},
			resp.BulkStringRespValue{
				Str: "apple",
			},
		},
	}

	setCommand, _ := cp.Parse(setResp)
	expected := SetCommand{
		Key:    "orange",
		Value:  "apple",
		Expiry: nil,
	}

	if setCommand != expected {
		t.Errorf("Should be set command")
	}
}

func TestSetPxCommand(t *testing.T) {
	setPxResp := resp.ArrayRespValue{
		Arr: []resp.RespValue{
			resp.BulkStringRespValue{
				Str: "SET",
			},
			resp.BulkStringRespValue{
				Str: "orange",
			},
			resp.BulkStringRespValue{
				Str: "apple",
			},
			resp.BulkStringRespValue{
				Str: "PX",
			},
			resp.BulkStringRespValue{
				Str: "100",
			},
		},
	}

	parsedCommand, _ := cp.Parse(setPxResp)
	currentTime := time.Now()

	expected := SetCommand{
		Key:    "orange",
		Value:  "apple",
		Expiry: &currentTime,
	}

	switch setPxCommand := parsedCommand.(type) {
	case SetCommand:
		if setPxCommand.Key != expected.Key {
			t.Errorf("Parsed command key should be: %s", expected.Key)
		}

		if setPxCommand.Value != expected.Value {
			t.Errorf("Parsed command value should be: %s", expected.Value)
		}

		difference := setPxCommand.Expiry.Sub(*expected.Expiry).Milliseconds()

		if difference < 95 || difference > 104 {
			t.Errorf("Expiry times dont match, difference does not fit threshold: %d", difference)
		}
	default:
		t.Errorf("Expected type: SetCommand")
	}
}

func TestGetCommand(t *testing.T) {
	getResp := resp.ArrayRespValue{
		Arr: []resp.RespValue{
			resp.BulkStringRespValue{
				Str: "GET",
			},
			resp.BulkStringRespValue{
				Str: "orange",
			},
		},
	}

	getCommand, _ := cp.Parse(getResp)
	expected := GetCommand{
		Key: "orange",
	}

	if getCommand != expected {
		t.Errorf("Should be set command")
	}
}
