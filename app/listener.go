package main

import (
	"fmt"
	"net"
)

func startListening() {
	listener, err := setUpListener()
	if err != nil {
		handleError("Failed to set up listener", err)
	}
	defer listener.Close()
	handleConnections(listener)
}

func setUpListener() (net.Listener, error) {
	listener, err := net.Listen("tcp", BindAddress)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func handleConnections(listener net.Listener) {
	fmt.Printf("Listening on http://%s\n", BindAddress)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn)
	}
}
