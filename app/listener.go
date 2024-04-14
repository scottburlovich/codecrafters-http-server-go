package main

import (
	"fmt"
	"net"
	"path/filepath"
)

func startListening(filesDir string) {
	listener, err := setUpListener()
	if err != nil {
		handleError("Failed to set up listener", err)
	}
	defer listener.Close()
	handleConnections(listener, filesDir)
}

func setUpListener() (net.Listener, error) {
	listener, err := net.Listen("tcp", BindAddress)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func handleConnections(listener net.Listener, filesDir string) {
	var err error
	fmt.Printf("Listening on http://%s\n", BindAddress)
	if filesDir != "" {
		filesDir, err = filepath.Abs(filesDir)
		if err != nil {
			fmt.Println("Error resolving absolute path of --directory argument: ", err)
			return
		}
		fmt.Printf("NOTE: /files route serving content from: %s\n", filesDir)
	}

	fmt.Println("Ready to accept connections, access logs will appear below...")
	fmt.Println("================================================================================")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn, filesDir)
	}
}
