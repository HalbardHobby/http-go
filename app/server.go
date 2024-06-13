package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	conn.Write(statusLine(200, "OK"))
	conn.Close()
}

func statusLine(statusCode int, message string) []byte {
	status := "HTTP/1.1" + strconv.Itoa(statusCode) + " " + message + "\r\n\r\n"
	return []byte(status)
}
