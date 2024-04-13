package main

import "strings"

type Route struct {
	Method string
	Path   string
}

type HandlerFunc func(params map[string]string) string

var router = map[Route]HandlerFunc{
	{Method: "GET", Path: "/"}:            rootHandler,
	{Method: "GET", Path: "/echo/:value"}: echoHandler,
}

func findRoute(req *HttpRequest) (HandlerFunc, map[string]string) {
	for route, handler := range router {
		if req.Method != route.Method {
			continue
		}

		routeSegments := strings.Split(route.Path, "/")
		pathSegments := strings.Split(req.URI, "/")

		if len(pathSegments) < len(routeSegments) {
			continue
		}

		params := make(map[string]string)
		matches := true

		for i, segment := range routeSegments {
			if strings.HasPrefix(segment, ":") {
				paramKey := strings.TrimPrefix(segment, ":")
				params[paramKey] = strings.Join(pathSegments[i:], "/")
				break
			} else if segment != pathSegments[i] {
				matches = false
				break
			}
		}

		if matches {
			return handler, params
		}
	}

	return notFoundHandler, nil
}
