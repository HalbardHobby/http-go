package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
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
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	ln, _, err := reader.ReadLine()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading request", err)
	}

	_, url := parseRequest(string(ln))

	if url == "/" {
		writer.WriteString(statusLine(200, "OK"))
	} else {
		writer.WriteString(statusLine(404, "Not Found"))
	}
	writer.Flush()
	conn.Close()
}

func parseRequest(req string) (method string, url string) {
	requestFields := strings.Split(req, " ")
	method = requestFields[0]
	url = requestFields[1]
	return
}

func statusLine(statusCode int, message string) string {
	return "HTTP/1.1 " + strconv.Itoa(statusCode) + " " + message + "\r\n\r\n"
}
