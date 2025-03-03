package main

import (
	"fmt"
	"net"
)

func main() {
	logs := NewLogger()
	logs.logInfo.Printf("Starting server on port :6379")

	// create a new server
	lsn, err := net.Listen("tcp", ":6379")
	if err != nil {
		logs.ServerError(err)
		return
	}
	// listen for connections
	con, err := lsn.Accept()
	if err != nil {
		logs.ServerError(err)
		return
	}
	defer func(con net.Conn) { _ = con.Close() }(con)
	for {
		resp := NewResp(con)
		value, err := resp.Read()
		if err != nil {
			logs.ServerError(err)
			return
		}
		fmt.Println(value)
		// ignore request and send back a PONG
		_, _ = con.Write([]byte("+OK\r\n"))
	}
}
