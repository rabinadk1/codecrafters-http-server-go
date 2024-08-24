package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Accepted connection for %s\n", conn.RemoteAddr())

	// Buffer to read from the connection
	buffer := make([]byte, 64)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading from connection: %s\n", err)
		return
	}

	headerParts := strings.SplitN(string(buffer[:n]), " ", 3)

	if len(headerParts) != 3 {
		fmt.Printf("Invalid request: %s\n", headerParts)
		return
	}

	requestTarget := headerParts[1]

	var response string
	if requestTarget == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n "
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("Error writing response: %s\n", err)
		return
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Fatalln("Failed to bind to port 4221")
	}
	defer l.Close()

	fmt.Println("Listening on port 4221")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn)

	}
}
