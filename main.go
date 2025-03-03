package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	logError := log.New(os.Stderr, "ERROR::\t", log.Ldate|log.Ltime|log.Lshortfile)
	logInfo := log.New(os.Stdout, "INFO::\t", log.Ldate|log.Ltime)

	logInfo.Printf("Starting server on port :6379")

	// create a new server
	lsn, err := net.Listen("tcp", ":6379")
	if err != nil {
		logError.Fatal(err)
		return
	}
	// listen for connections
	con, err := lsn.Accept()
	if err != nil {
		logError.Fatal(err)
		return
	}
	defer func(con net.Conn) { _ = con.Close() }(con)
	for {
		rsp := &resp{reader: bufio.NewReader(con)}
		val, err := rsp.Read()
		if err != nil {
			logError.Fatal(err)
			return
		}
		fmt.Println(val)
		// ignore request and send back a PONG
		_, _ = con.Write([]byte("+OK\r\n"))
	}
}
