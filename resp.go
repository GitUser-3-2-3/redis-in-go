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

func (r *Resp) readLine() (line []byte, bytlen int, err error) {
	for {
		byt, err := r.reader.ReadByte()
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

func (r *Resp) readInteger() (x, bytlen int, err error) {
	line, bytlen, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, bytlen, err
	}
	return int(i64), bytlen, nil
}

func (r *Resp) Read() (Value, error) {
	bytTyp, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch bytTyp {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(bytTyp))
		return Value{}, nil
	}
}

func (r *Resp) readArray() (Value, error) {
	val := Value{}
	val.typ = "array"

	ln, _, err := r.readInteger()
	if err != nil {
		return val, err
	}
	val.array = make([]Value, ln)

	for i := 0; i < ln; i++ {
		v, err := r.Read()
		if err != nil {
			return val, err
		}
		val.array[i] = v
	}
	return val, nil
}

func (r *Resp) readBulk() (Value, error) {
	val := Value{}
	val.typ = "bulk"

	ln, _, err := r.readInteger()
	if err != nil {
		return val, err
	}
	bulk := make([]byte, ln)
	_, _ = r.reader.Read(bulk)

	val.bulk = string(bulk)

	// read the trailing CRLF
	_, _, _ = r.readLine()

	return val, nil
}
