package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var directory *string

func main() {

	directory = flag.String("directory", "/tmp", "Location for the files to be served.")
	flag.Parse()
	fmt.Fprintf(os.Stdout, "Serving files from: %s\n", *directory)

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

	_, path := parseRequest(string(ln))
	splitPath := strings.Split(path, "/")
	if path == "/" {
		conn.Write([]byte(statusLine(200, "OK") + "\r\n"))
	} else if splitPath[1] == "echo" {
		conn.Write([]byte(statusLine(200, "OK")))
		msg := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(splitPath[2]), splitPath[2])
		conn.Write([]byte(msg))
	} else if splitPath[1] == "user-agent" {
		headers := getHeaders(reader)
		agent := headers["User-Agent"]
		conn.Write([]byte(statusLine(200, "OK")))
		msg := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(agent), agent)
		conn.Write([]byte(msg))
	} else if splitPath[1] == "files" {
		target := *directory + splitPath[2]
		_, err := os.Stat(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "problem getting path: %s", err.Error())
			conn.Write([]byte(statusLine(404, "Not Found") + "\r\n"))
		} else {
			file, _ := os.Open(target)
			info, _ := file.Stat()
			fileContent := make([]byte, info.Size())
			file.Read(fileContent)
			msg := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(fileContent), fileContent)
			conn.Write([]byte(msg))
		}
	} else {
		conn.Write([]byte(statusLine(404, "Not Found") + "\r\n"))
	}
}

func getHeaders(reader *bufio.Reader) map[string]string {
	resp := make(map[string]string)

	ln, _, _ := reader.ReadLine()

	for !bytes.Equal(ln, []byte("")) {
		fields := strings.Split(string(ln), ": ")
		resp[fields[0]] = fields[1]
		ln, _, _ = reader.ReadLine()
	}
	return resp
}

func parseRequest(req string) (method string, url string) {
	requestFields := strings.Split(req, " ")
	method = requestFields[0]
	url = requestFields[1]
	return
}

func statusLine(statusCode int, message string) string {
	return "HTTP/1.1 " + strconv.Itoa(statusCode) + " " + message + "\r\n"
}
