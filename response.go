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
