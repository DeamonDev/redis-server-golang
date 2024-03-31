package resp

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

var rp *RespParser

func TestMain(m *testing.M) {
	rp = NewParser()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestParseString(t *testing.T) {
	str := "+HELLO_RESP_PARSER\r\n"
	byteSlice := []byte(str)

	reader := bytes.NewReader(byteSlice)

	expectedString := "HELLO_RESP_PARSER"
	expectedRespValue := StringRespValue{
		Str: "HELLO_RESP_PARSER",
	}

	parsedValue, _ := rp.Parse(reader)

	if parsedValue != expectedRespValue {
		t.Errorf("Expected StringRespValue: %s", expectedString)
	}
}

func TestParseNumber(t *testing.T) {
	str := ":50\r\n"
	byteSlice := []byte(str)

	reader := bytes.NewReader(byteSlice)

	exptectedNum := 50
	expectedRespValue := NumberRespValue{
		Num: exptectedNum,
	}

	parsedValue, _ := rp.Parse(reader)

	if parsedValue != expectedRespValue {
		t.Errorf("Expected NumberRespValue: %d", exptectedNum)
	}
}

func TestParseBulkString(t *testing.T) {
	str := "$5\r\nHELLO\r\n"
	byteSlice := []byte(str)

	reader := bytes.NewReader(byteSlice)

	expectedBulkString := "HELLO"
	expectedRespValue := BulkStringRespValue{
		Str: expectedBulkString,
	}

	parsedValue, _ := rp.Parse(reader)

	if parsedValue != expectedRespValue {
		t.Errorf("Expected BulkStringRespValue: %s", expectedBulkString)
	}
}

func TestParseArray(t *testing.T) {
	str := "*3\r\n+ARR_1\r\n:20\r\n$5\r\nHELLO\r\n"
	byteSlice := []byte(str)

	reader := bytes.NewReader(byteSlice)

	expectedRespValue := ArrayRespValue{
		Arr: []RespValue{
			StringRespValue{
				Str: "ARR_1",
			},
			NumberRespValue{
				Num: 20,
			},
			BulkStringRespValue{
				Str: "HELLO",
			},
		},
	}

	parsedValue, _ := rp.Parse(reader)

	switch v := parsedValue.(type) {
	case ArrayRespValue:
		if !cmp.Equal(v, expectedRespValue) {
			t.Errorf("NOT EQUAL")
		}
	default:
		t.Errorf("WRONG")
	}
}
