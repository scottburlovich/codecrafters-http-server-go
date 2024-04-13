package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	BindAddress    = "0.0.0.0:4221"
	HttpOKResponse = "HTTP/1.1 200 OK\r\n\r\n"
)

type HttpRequest struct {
	Method          string
	URI             string
	ProtocolVersion string
	Headers         map[string]string
}

func main() {
	fmt.Println("Logs from your program will appear here!")
	listener, err := net.Listen("tcp", BindAddress)
	if err != nil {
		handleError("Failed to bind to port 4221", err)
	}

	conn, err := listener.Accept()
	if err != nil {
		handleError("Error accepting connection", err)
	}

	handleConnection(conn)
}

func handleError(message string, err error) {
	fmt.Println(message+": ", err.Error())
	os.Exit(1)
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		handleError("Error reading", err)
	}
	req, err := parseHttpRequest(buf)
	if err != nil {
		handleError("Error parsing HTTP request", err)
	}

	handleHttpResponse(conn, req)
}

func handleHttpResponse(conn net.Conn, req *HttpRequest) {
	var err error
	var response string
	if req.URI == "/" {
		response = HttpOKResponse
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		handleError("Error writing", err)
	}
}

func parseHttpRequest(request []byte) (*HttpRequest, error) {
	lines := strings.Split(string(request), "\r\n")

	if len(lines) < 1 {
		return nil, errors.New("invalid HTTP request")
	}

	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) != 3 {
		return nil, errors.New("invalid request line")
	}

	httpRequest := &HttpRequest{
		Method:          requestLine[0],
		URI:             requestLine[1],
		ProtocolVersion: requestLine[2],
		Headers:         make(map[string]string),
	}

	for _, line := range lines[1:] {
		if line == "" {
			break
		}
		header := strings.SplitN(line, ": ", 2)
		if len(header) != 2 {
			return nil, errors.New("invalid header: " + line)
		}
		httpRequest.Headers[header[0]] = header[1]
	}

	return httpRequest, nil
}
