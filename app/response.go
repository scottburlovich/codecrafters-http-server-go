package main

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func buildResponseHeaders(contentType string, contentLength string, customHeaders *map[string]string) string {
	headers := []string{
		"Content-Type: " + contentType,
		"Content-Length: " + contentLength,
	}

	if customHeaders != nil {
		for key, value := range *customHeaders {
			customHeader := key + ": " + value
			headers = append(headers, customHeader)
		}
	}

	return strings.Join(headers, Separator) + Separator
}

func buildResponse(args HandlerArgs, status, body string, customContentType *string, customHeaders *map[string]string) {
	var contentType string
	if customContentType != nil {
		contentType = *customContentType
	} else {
		contentType = detectContentType([]byte(body))
	}
	contentLength := strconv.Itoa(len(body))

	headers := buildResponseHeaders(contentType, contentLength, customHeaders)
	var response string
	if headers == "" && body == "" {
		response = status + Separator + Separator
	}
	response = status + headers + Separator + body

	handleResponse(args.Conn, response)
	outputAccessLog(args, status, contentLength)
}

func outputAccessLog(args HandlerArgs, status string, contentLength string) {
	currentTime := time.Now().UTC()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	re := regexp.MustCompile(`\d{3}`)
	code := re.FindString(status)
	requestStr := fmt.Sprintf("%s %s %s", args.Req.Method, args.Req.URI, args.Req.ProtocolVersion)

	fmt.Printf(
		"%s - - [%s] \"%s\" %s %s \"%s\"\n",
		args.Req.ClientIP,
		formattedTime,
		requestStr,
		code,
		contentLength,
		args.Req.Headers["User-Agent"],
	)
}

func handleResponse(conn net.Conn, response string) {
	_, err := conn.Write([]byte(response))
	if err != nil {
		handleError("Error writing response to connection", err)
	}
}
