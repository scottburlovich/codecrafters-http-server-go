package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type HttpRequest struct {
	Method          string
	URI             string
	ClientIP        string
	ProtocolVersion string
	Headers         map[string]string
}

func handleConnection(conn net.Conn, filesDir string) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}

	req, err := parseRequest(buf)
	if err != nil {
		fmt.Println("Error parsing HTTP request: ", err.Error())
		return
	}

	// Get the client IP address from the connection
	req.ClientIP = conn.RemoteAddr().String()

	handleRequest(conn, req, filesDir)
}

func handleRequest(conn net.Conn, req *HttpRequest, filesDir string) {
	args := HandlerArgs{
		Conn:     conn,
		Req:      req,
		Params:   nil,
		FileRoot: filesDir,
	}

	handler, args := findRoute(args)
	handler(args)
}

func parseRequest(request []byte) (*HttpRequest, error) {
	lines := strings.Split(string(request), "\r\n")

	if len(lines) < 1 {
		return nil, errors.New("invalid request")
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
