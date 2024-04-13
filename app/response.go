package main

import (
	"net"
	"strings"
)

func buildResponse(status, headers, body string) string {
	if headers == "" && body == "" {
		return status + Separator + Separator
	}
	return status + headers + Separator + body
}

func buildResponseHeaders(contentType string, contentLength string) string {
	headers := []string{
		"Content-Type: " + contentType,
		"Content-Length: " + contentLength,
	}

	return strings.Join(headers, Separator) + Separator
}

func handleResponse(conn net.Conn, req *HttpRequest) {
	var err error
	response := handleRequest(req)

	_, err = conn.Write([]byte(response))
	if err != nil {
		handleError("Error writing response to connection", err)
	}
}
