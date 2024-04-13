package main

import "net/http"

func detectContentType(content []byte) string {
	inferred := http.DetectContentType(content)
	if inferred == "text/plain; charset=utf-8" {
		return "text/plain"
	} else {
		return inferred
	}
}
