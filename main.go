package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	logError := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	logInfo := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

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
		buf := make([]byte, 1024)

		// read message from client
		_, err := con.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			logInfo.Println("Error reading from client: ", err)
			os.Exit(1)
		}
		// ignore request and send back a PONG
		_, _ = con.Write([]byte("+OK\r\n"))
	}
}
