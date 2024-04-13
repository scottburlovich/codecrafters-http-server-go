package main

import (
	"net"
	"strconv"
	"strings"
)

func buildResponseHeaders(contentType string, contentLength string, customHeaders *map[string]string) string {
	headers := []string{
		"Content-Type: " + contentType,
		"Content-Length: " + contentLength,
	}

	if customHeaders != nil { // check if customHeaders is not nil
		for key, value := range *customHeaders {
			customHeader := key + ": " + value
			headers = append(headers, customHeader)
		}
	}

	return strings.Join(headers, Separator) + Separator
}

func buildResponse(status, body string, customHeaders *map[string]string) string {
	contentType := detectContentType([]byte(body))
	contentLength := strconv.Itoa(len(body))

	headers := buildResponseHeaders(contentType, contentLength, customHeaders)

	if headers == "" && body == "" {
		return status + Separator + Separator
	}
	return status + headers + Separator + body
}

func handleResponse(conn net.Conn, req *HttpRequest) {
	var err error
	response := handleRequest(req)

	_, err = conn.Write([]byte(response))
	if err != nil {
		handleError("Error writing response to connection", err)
	}
}
