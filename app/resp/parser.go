package resp

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type RespParserError struct {
	Domain   string
	Message  string
	Internal error
}

func (e RespParserError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Domain, e.Message, e.Internal)
}

func NewRespParserError(message string, err error) *RespParserError {
	return &RespParserError{
		Domain:   "RespParser",
		Message:  message,
		Internal: err,
	}
}

type RespValue interface {
	RespValueType() string
}

type StringRespValue struct {
	Str string
}

func (_ StringRespValue) RespValueType() string {
	return "STRING_RESP_VALUE"
}

type NumberRespValue struct {
	Num int
}

func (_ NumberRespValue) RespValueType() string {
	return "NUMBER_RESP_VALUE"
}

type BulkStringRespValue struct {
	Str string
}

func (_ BulkStringRespValue) RespValueType() string {
	return "BULK_STRING_RESP_VALUE"
}

type ArrayRespValue struct {
	Arr []RespValue
}

func (_ ArrayRespValue) RespValueType() string {
	return "ARRAY_RESP_VALUE"
}

const (
	STRING = '+'
	NUMBER = ':'
	BULK   = '$'
	ARRAY  = '*'
)

type RespParser struct {
}

func NewParser() *RespParser {
	return &RespParser{}
}

func (respParser *RespParser) Parse(reader *bytes.Reader) (RespValue, error) {
	firstByte, _ := reader.ReadByte()

	switch firstByte {
	case STRING:
		return respParser.parseString(reader)
	case BULK:
		return respParser.parseBulkString(reader)
	case NUMBER:
		return respParser.parseNumber(reader)
	case ARRAY:
		return respParser.parseArray(reader)
	default:
		return nil, NewRespParserError("first byte is unknown", nil)
	}
}

func (respParser *RespParser) parseString(reader *bytes.Reader) (RespValue, error) {
	var textBuffer bytes.Buffer

	err := readRespLine(reader, &textBuffer)
	if err != nil {
		return nil, err
	}

	text := textBuffer.String()
	stringRespValue := StringRespValue{
		Str: text,
	}

	return stringRespValue, nil
}

func (respParser *RespParser) parseNumber(reader *bytes.Reader) (RespValue, error) {
	var textBuffer bytes.Buffer

	err := readRespLine(reader, &textBuffer)
	if err != nil {
		return nil, err
	}

	text := textBuffer.String()
	num, _ := strconv.Atoi(text)
	numberRespValue := NumberRespValue{
		Num: num,
	}

	return numberRespValue, nil
}

// happy path af
func (respParser *RespParser) parseBulkString(reader *bytes.Reader) (RespValue, error) {
	length, err := readIntCRLF(reader)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)

	io.ReadFull(reader, data)

	reader.ReadByte()
	reader.ReadByte()

	bulkStringRespValue := BulkStringRespValue{
		Str: string(data),
	}

	return bulkStringRespValue, nil
}

func readIntCRLF(reader *bytes.Reader) (int, error) {
	var textBuffer bytes.Buffer

	for {
		if b, err := reader.ReadByte(); err != nil {
			return 0, err
		} else if b == '\r' {
			break
		} else {
			textBuffer.WriteByte(b)
		}
	}

	if b, err := reader.ReadByte(); err != nil {
		return 0, err
	} else if b != '\n' {
		return 0, fmt.Errorf("invalid RESP format: expected '\\n', got '%c'", b)
	}

	var length int

	_, err := fmt.Fscanf(&textBuffer, "%d", &length)
	if err != nil {
		return 0, NewRespParserError("error while reading CRLF int", err)
	}

	return length, nil
}

func (respParser *RespParser) parseArray(reader *bytes.Reader) (RespValue, error) {
	length, err := readIntCRLF(reader)
	if err != nil {
		return nil, err
	}

	arr := make([]RespValue, length)

	for i := 0; i < length; i++ {
		respValue, _ := respParser.Parse(reader)
		arr[i] = respValue
	}

	arrayRespValue := ArrayRespValue{
		Arr: arr,
	}

	return arrayRespValue, nil
}

func readRespLine(reader *bytes.Reader, textBuffer *bytes.Buffer) error {
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return NewRespParserError("error while reading first byte", err)
		}

		if b == '\r' {
			c, _ := reader.ReadByte()
			if c == '\n' {
				break
			}

			return NewRespParserError("There has to be newline character after CR", nil)
		}

		textBuffer.WriteByte(b)
	}
	return nil
}
