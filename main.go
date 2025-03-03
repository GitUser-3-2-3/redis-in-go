package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	log := initLogs()
	log.logInfo.Printf("Starting server on port :6379")

	// create a new server
	lsn, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.serverError(err)
		return
	}
	// listen for connections
	con, err := lsn.Accept()
	if err != nil {
		log.serverError(err)
		return
	}
	defer func(con net.Conn) { _ = con.Close() }(con)
	for {
		rsp := &resp{reader: bufio.NewReader(con)}
		val, err := rsp.Read()
		if err != nil {
			log.serverError(err)
			return
		}
		fmt.Println(val)
		// ignore request and send back a PONG
		_, _ = con.Write([]byte("+OK\r\n"))
	}
}
