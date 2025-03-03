package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// todo -> better error handling

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

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{bufio.NewReader(rd)}
}

func (resp *Resp) readLine() (line []byte, bytlen int, err error) {
	for {
		byt, err := resp.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		bytlen += 1
		line = append(line, byt)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], bytlen, nil
}

func (resp *Resp) readInteger() (x, bytlen int, err error) {
	line, bytlen, err := resp.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, bytlen, err
	}
	return int(i64), bytlen, nil
}

func (resp *Resp) Read() (Value, error) {
	bytTyp, err := resp.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch bytTyp {
	case ARRAY:
		return resp.readArray()
	case BULK:
		return resp.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(bytTyp))
		return Value{}, nil
	}
}

func (resp *Resp) readArray() (Value, error) {
	val := Value{}
	val.typ = "array"

	ln, _, err := resp.readInteger()
	if err != nil {
		return val, err
	}
	val.array = make([]Value, ln)

	for i := 0; i < ln; i++ {
		v, err := resp.Read()
		if err != nil {
			return val, err
		}
		val.array[i] = v
	}
	return val, nil
}

func (resp *Resp) readBulk() (Value, error) {
	val := Value{}
	val.typ = "bulk"

	ln, _, err := resp.readInteger()
	if err != nil {
		return val, err
	}
	bulk := make([]byte, ln)
	_, _ = resp.reader.Read(bulk)

	val.bulk = string(bulk)

	// read the trailing CRLF
	_, _, _ = resp.readLine()

	return val, nil
}
