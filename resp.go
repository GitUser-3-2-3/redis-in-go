package main

import (
	"bufio"
	"fmt"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []value
}

type resp struct {
	reader *bufio.Reader
}

func (rsp *resp) readLine() (line []byte, bytlen int, err error) {
	for {
		byt, err := rsp.reader.ReadByte()
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

func (rsp *resp) readInteger() (x, bytlen int, err error) {
	line, bytlen, err := rsp.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, bytlen, err
	}
	return int(i64), bytlen, nil
}

func (rsp *resp) Read() (value, error) {
	bytTyp, err := rsp.reader.ReadByte()
	if err != nil {
		return value{}, err
	}
	switch bytTyp {
	case ARRAY:
		return rsp.readArray()
	case BULK:
		return rsp.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(bytTyp))
		return value{}, nil
	}
}

func (rsp *resp) readArray() (value, error) {
	val := value{}
	val.typ = "array"

	ln, _, err := rsp.readInteger()
	if err != nil {
		return val, err
	}
	val.array = make([]value, ln)

	for i := 0; i < ln; i++ {
		v, err := rsp.Read()
		if err != nil {
			return val, err
		}
		val.array = append(val.array, v)
	}
	return val, nil
}

func (rsp *resp) readBulk() (value, error) {
	val := value{}
	val.typ = "bulk"

	ln, _, err := rsp.readInteger()
	if err != nil {
		return val, err
	}
	bulk := make([]byte, ln)

	_, err = rsp.reader.Read(bulk)
	if err != nil {
		return value{}, err
	}
	val.bulk = string(bulk)

	// read the trailing CRLF
	_, _, err = rsp.readLine()
	if err != nil {
		return value{}, err
	}
	return val, nil
}
