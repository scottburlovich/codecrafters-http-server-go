package main

import (
	"net"
	"os"
	"path/filepath"
)

type HandlerArgs struct {
	Conn     net.Conn
	Req      *HttpRequest
	Params   map[string]string
	FileRoot string
}

func echoHandler(args HandlerArgs) {
	echoString := args.Params["value"]
	buildResponse(args, HttpOKResponse, echoString, nil, nil)
}

func fileDownloadHandler(args HandlerArgs) {
	if args.FileRoot == "" {
		buildResponse(args, HttpInternalErrorResponse, "No directory set", nil, nil)
	}

	filename, ok := args.Params["filename"]
	if !ok {
		buildResponse(args, HttpNotFoundResponse, "", nil, nil)
	}

	filePath := filepath.Join(args.FileRoot, filename)
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		buildResponse(args, HttpNotFoundResponse, "", nil, nil)
	} else {
		responseContentType := "application/octet-stream"
		buildResponse(args, HttpOKResponse, string(fileContent), &responseContentType, nil)
	}
}

func fileUploadHandler(args HandlerArgs) {
	if args.FileRoot == "" {
		buildResponse(args, HttpInternalErrorResponse, "No directory set", nil, nil)
	}

	filename, ok := args.Params["filename"]
	if !ok {
		buildResponse(args, HttpInternalErrorResponse, "Filename required", nil, nil)
	}

	filePath := filepath.Join(args.FileRoot, filename)
	err := os.WriteFile(filePath, []byte(args.Req.Body), 0644)
	if err != nil {
		buildResponse(args, HttpInternalErrorResponse, "", nil, nil)
	} else {
		buildResponse(args, HttpCreatedResponse, "", nil, nil)
	}

}

func rootHandler(args HandlerArgs) {
	buildResponse(args, HttpOKResponse, "", nil, nil)
}

func notFoundHandler(args HandlerArgs) {
	buildResponse(args, HttpNotFoundResponse, "", nil, nil)
}

func userAgentHandler(args HandlerArgs) {
	userAgent := args.Req.Headers["User-Agent"]
	buildResponse(args, HttpOKResponse, userAgent, nil, nil)
}
