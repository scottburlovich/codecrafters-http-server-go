package main

import "strings"

type Route struct {
	Method string
	Path   string
}

type HandlerFunc func(args HandlerArgs)

func getRouter(args HandlerArgs) map[Route]HandlerFunc {
	return map[Route]HandlerFunc{
		{Method: "GET", Path: "/"}:                rootHandler,
		{Method: "GET", Path: "/echo/:value"}:     echoHandler,
		{Method: "GET", Path: "/user-agent"}:      userAgentHandler,
		{Method: "GET", Path: "/files/:filename"}: filesHandler,
	}
}

func findRoute(args HandlerArgs) (HandlerFunc, HandlerArgs) {
	router := getRouter(args)
	for route, handler := range router {
		if args.Req.Method != route.Method {
			continue
		}

		routeSegments := strings.Split(route.Path, "/")
		pathSegments := strings.Split(args.Req.URI, "/")

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

		// Populate the args Params field
		args.Params = params

		if matches {
			return handler, args
		}
	}

	return notFoundHandler, args
}
