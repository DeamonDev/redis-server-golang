package resp

import (
	"bytes"
	"errors"
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

func TestReturnsErrorWhenProvidedIncorrectString(t *testing.T) {
	str := "INCORRECT_RESP_STRING"
	byteSlice := []byte(str)

	reader := bytes.NewReader(byteSlice)

	expectedRespError := NewRespParserError("first byte is unknown", nil)

	_, err := rp.Parse(reader)

	var actualRespError *RespParserError
	ok := errors.As(err, &actualRespError)
	if !ok {
		t.Errorf("expected %v to be of type RespParserError", err)
		return
	}

	if actualRespError.Message != expectedRespError.Message {
		t.Errorf("expected %v, actual %v", expectedRespError, actualRespError)
	}
}

func TestReturnsErrorWhenProvidedStringWithoutCRTokenAtTheEnd(t *testing.T) {
	str := "+TEXT\r"
	byteSlice := []byte(str)

	reader := bytes.NewReader(byteSlice)

	expectedRespError := NewRespParserError("There has to be newline character after CR", nil)

	_, err := rp.Parse(reader)

	var actualRespError *RespParserError
	ok := errors.As(err, &actualRespError)
	if !ok {
		t.Errorf("expected %v to be of type RespParserError", err)
		return
	}

	if actualRespError.Message != expectedRespError.Message {
		t.Errorf("expected %v, actual %v", expectedRespError, actualRespError)
	}

}

func TestReturnsErrorWhenProvidedNumberWithoutCRTokenAtTheEnd(t *testing.T) {
	str := ":42\r"
	byteSlice := []byte(str)

	reader := bytes.NewReader(byteSlice)

	expectedRespError := NewRespParserError("There has to be newline character after CR", nil)

	_, err := rp.Parse(reader)

	var actualRespError *RespParserError
	ok := errors.As(err, &actualRespError)
	if !ok {
		t.Errorf("expected %v to be of type RespParserError", err)
		return
	}

	if actualRespError.Message != expectedRespError.Message {
		t.Errorf("expected %v, actual %v", expectedRespError, actualRespError)
	}

}

func TestReturnsErrorWhenProvidedIncorrectIntCRLF(t *testing.T) {
	str := "$xx\r\nhh\r\n"
	byteSlice := []byte(str)

	reader := bytes.NewReader(byteSlice)

	expectedRespError := NewRespParserError("error while reading CRLF int", errors.New("expected integer"))

	_, err := rp.Parse(reader)

	var actualRespError *RespParserError
	ok := errors.As(err, &actualRespError)
	if !ok {
		t.Errorf("expected %v to be of type RespParserError", err)
		return
	}

	if actualRespError.Message != expectedRespError.Message {
		t.Errorf("expected %v, actual %v", expectedRespError, actualRespError)
	}

	if actualRespError.Internal.Error() != expectedRespError.Internal.Error() {
		t.Errorf("expected %s, actual %s", expectedRespError.Internal.Error(), actualRespError.Internal.Error())
	}

}
