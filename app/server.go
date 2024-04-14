package main

import (
	"flag"
	"fmt"
)

var dirPtr = flag.String("directory", "", "directory to serve files from for the /files route")

func main() {
	flag.Parse()
	fmt.Println("HTTP server initializing...")
	startListening(*dirPtr)
}
