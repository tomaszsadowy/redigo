package main

import (
	"bufio"
	"io"
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
