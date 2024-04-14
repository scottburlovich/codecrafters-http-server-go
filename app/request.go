package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"net"
	"strconv"
	"strings"
)

type HttpRequest struct {
	Method          string
	URI             string
	ClientIP        string
	ProtocolVersion string
	Headers         map[string]string
	Body            string
	BodyContentType string
}

func handleConnection(conn net.Conn, filesDir string) {
	defer conn.Close()

	var buf bytes.Buffer
	reader := bufio.NewReader(conn)
	contentLength := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			handleError("Failed to read request line", err)
			return
		}
		line = strings.TrimSuffix(line, "\r\n")
		buf.WriteString(line + "\r\n")

		if line == "" {
			break
		}

		if strings.HasPrefix(line, "Content-Length:") {
			size := strings.Trim(strings.TrimPrefix(line, "Content-Length:"), " ")
			contentLength, err = strconv.Atoi(size)
			if err != nil {
				handleError("Failed to parse Content-Length", err)
				return
			}
		}
	}

	if contentLength > 0 {
		body := make([]byte, contentLength)
		bytesRead, err := reader.Read(body)
		if err != nil || bytesRead != contentLength {
			handleError("Failed to read request body", err)
			return
		}
		buf.Write(body)
	}

	request, err := parseRequest(buf.Bytes())
	if err != nil {
		handleError("Failed to parse request", err)
		return
	}

	request.ClientIP = conn.RemoteAddr().String()
	handleRequest(conn, request, filesDir)
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

	bodyLinesIndex := 0
	for i, line := range lines[1:] {
		if line == "" {
			bodyLinesIndex = i + 2
			break
		}
		header := strings.SplitN(line, ": ", 2)
		if len(header) != 2 {
			return nil, errors.New("invalid header: " + line)
		}
		httpRequest.Headers[header[0]] = header[1]
	}
	httpRequest.Body = strings.Join(lines[bodyLinesIndex:], "\r\n")

	if httpRequest.Method == "POST" && strings.Contains(httpRequest.Headers["Content-Type"], "multipart/form-data") {
		mediaType, params, err := mime.ParseMediaType(httpRequest.Headers["Content-Type"])
		if err != nil {
			return nil, errors.New("invalid Content-Type header: " + err.Error())
		}
		httpRequest.BodyContentType = mediaType

		if mediaType == "multipart/form-data" {
			body := strings.NewReader(httpRequest.Body)
			mr := multipart.NewReader(body, params["boundary"])

			formData, err := mr.ReadForm(-1)
			if err != nil {
				return nil, errors.New("invalid multipart form: " + err.Error())
			}

			fmt.Println("Form data:", formData.Value)
			fmt.Println("Files:", formData.File)
		}
	}

	return httpRequest, nil
}
