package main

func echoHandler(req *HttpRequest, params map[string]string) string {
	echoString := params["value"]
	return buildResponse(HttpOKResponse, echoString, nil)
}

func rootHandler(req *HttpRequest, params map[string]string) string {
	return buildResponse(HttpOKResponse, "", nil)
}

func notFoundHandler(req *HttpRequest, params map[string]string) string {
	return buildResponse(HttpNotFoundResponse, "", nil)
}

func userAgentHandler(req *HttpRequest, params map[string]string) string {
	userAgent := req.Headers["User-Agent"]
	return buildResponse(HttpOKResponse, userAgent, nil)
}
