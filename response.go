package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Response struct {
	reader *bufio.Reader
}

func NewResponse(rd io.Reader) *Response {
	return &Response{reader: bufio.NewReader(rd)}
}

func (r *Response) readLine() (line []byte, n int, err error) {
	for {
		_byte, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, _byte)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Response) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Response) readArray() (Value, error) {
	val := Value{}
	val.typ = "array"

	len, _, err := r.readInteger()
	if err != nil {
		return val, err
	}

	val.array = make([]Value, 0)
	for i := 0; i < len; i++ {
		_val, err := r.Read()
		if err != nil {
			return val, err
		}

		val.array = append(val.array, _val)
	}
	return val, nil
}

func (r *Response) readBulk() (Value, error) {
	val := Value{}
	val.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return val, err
	}

	bulk := make([]byte, len)
	r.reader.Read(bulk)
	val.bulk = string(bulk)
	r.readLine()

	return val, nil
}

func (r *Response) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

func (val Value) Marshal() []byte {
	switch val.typ {
	case "array":
		return val.marshalArray()
	case "bulk":
		return val.marshalBulk()
	case "string":
		return val.marshalString()
	case "null":
		return val.marshallNull()
	case "error":
		return val.marshallError()
	default:
		return []byte{}
	}
}

func (val Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, val.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (val Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(val.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, val.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (val Value) marshalArray() []byte {
	len := len(val.array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, val.array[i].Marshal()...)
	}

	return bytes
}

func (val Value) marshallError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, val.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (val Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}
