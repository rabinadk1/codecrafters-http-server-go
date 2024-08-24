package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn, directory string) {
	defer conn.Close()

	fmt.Printf("Accepted connection for %s\n", conn.RemoteAddr())

	// Buffer to read from the connection
	buffer := make([]byte, 1204)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading from connection: %s\n", err)
		return
	}
	buffer = buffer[:n]
	stringBuffer := string(buffer)

	headerParts := strings.SplitN(stringBuffer, " ", 3)

	if len(headerParts) != 3 {
		fmt.Printf("Invalid request: %s\n", headerParts)
		return
	}

	requestTarget := headerParts[1]

	response := "HTTP/1.1 404 Not Found\r\n\r\n"
	if requestTarget == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	} else if requestTarget == "/user-agent" {
		const userAgentPrefix = "User-Agent: "
		userAgentIndex := strings.Index(stringBuffer, userAgentPrefix)
		if userAgentIndex == -1 {
			fmt.Printf("Could not find User-Agent in request: %s\n", stringBuffer)
		}

		remString := stringBuffer[userAgentIndex+len(userAgentPrefix):]

		rnIndex := strings.Index(remString, "\r\n")
		if rnIndex == -1 {
			fmt.Printf("Could not find \\r\\n in request: %s\n", stringBuffer)
		}

		userAgent := remString[:rnIndex]
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)

	} else if prefix := "/echo/"; strings.HasPrefix(requestTarget, prefix) {
		echoBody := strings.TrimPrefix(requestTarget, prefix)
		contentLength := len(echoBody)
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, echoBody)
	} else if prefix := "/files/"; strings.HasPrefix(requestTarget, prefix) {
		fileName := strings.TrimPrefix(requestTarget, prefix)
		filePath := directory + "/" + fileName

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("File Not Found: %s", err)
		} else {
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(fileContent), string(fileContent))
		}
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("Error writing response: %s\n", err)
		return
	}
}

func main() {
	directory := "."
	if len(os.Args) > 1 && os.Args[1] == "--directory" {
		if len(os.Args) < 3 {
			log.Fatalln("Usage: ./http_server --directory <directory>")
		}
		directory = strings.TrimRight(os.Args[2], "/")
	}

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

		go handleConnection(conn, directory)

	}
}
