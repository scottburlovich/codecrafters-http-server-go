package main

import "fmt"

func echoHandler(params map[string]string) string {
	echoString := params["value"]
	responseHeaders := buildResponseHeaders("text/plain", fmt.Sprintf("%d", len(echoString)))
	return buildResponse(HttpOKResponse, responseHeaders, echoString)
}

func rootHandler(params map[string]string) string {
	return buildResponse(HttpOKResponse, "", "")
}

func notFoundHandler(params map[string]string) string {
	return buildResponse(HttpNotFoundResponse, "", "")
}
