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
	defer conn.Close()

	reader := bufio.NewReader(conn)
	ln, _, err := reader.ReadLine()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading request", err)
	}

	_, url := parseRequest(string(ln))
	if url == "/" {
		conn.Write([]byte(statusLine(200, "OK")))
	} else {
		conn.Write([]byte(statusLine(404, "Not Found")))
	}
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
