package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func handleUserAgent(stringBuffer string) (string, error) {
	const userAgentPrefix = "User-Agent: "
	userAgent, err := getHeaderValue(userAgentPrefix, stringBuffer)
	if err != nil {
		return "", fmt.Errorf("Error getting %s%s", userAgentPrefix, err)
	}

	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent), nil
}

func handleFiles(requestMethod string, stringBuffer string, filePath string) (string, error) {
	if requestMethod != "POST" {
		fileContent, err := os.ReadFile(filePath)

		if os.IsNotExist(err) {
			return "HTTP/1.1 404 Not Found\r\n\r\n", nil
		}

		if err != nil {
			return "", fmt.Errorf("File Not Found: %s", err)
		}

		return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(fileContent), string(fileContent)), nil
	}

	const contentTypePrefix = "Content-Length: "
	contentLengthStr, err := getHeaderValue(contentTypePrefix, stringBuffer)
	if err != nil {
		return "", fmt.Errorf("Error getting %s%s", contentLengthStr, err)
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		return "", fmt.Errorf("Could not convert content length to integer from request: %s\n", stringBuffer)
	}

	const headerBodySeperator = "\r\n\r\n"
	headerBodySeperatorIndex := strings.Index(stringBuffer, headerBodySeperator)
	if headerBodySeperatorIndex == -1 {
		return "", fmt.Errorf("Could not find \\r\\n\\r\\n in request: %s\n", stringBuffer)
	}

	body := stringBuffer[headerBodySeperatorIndex+len(headerBodySeperator):]

	if contentLength != len(body) {
		return "", fmt.Errorf("Content-Length (%d) did not match body length: %d\n", contentLength, len(body))
	}

	if err := os.WriteFile(filePath, []byte(body), 0644); err != nil {
		return "", fmt.Errorf("Could not write to file: %s\n", err)
	}

	return "HTTP/1.1 201 Created\r\n\r\n", nil
}

func handleConnection(conn net.Conn, directory string) {
	defer conn.Close()

	fmt.Printf("Accepted connection for %s\n", conn.RemoteAddr())

	// Buffer to read from the connection
	buffer := make([]byte, 1024)

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

	requestMethod, requestTarget := headerParts[0], headerParts[1]

	response := "HTTP/1.1 404 Not Found\r\n\r\n"
	if requestTarget == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	} else if requestTarget == "/user-agent" {
		response, err = handleUserAgent(stringBuffer)
		if err != nil {
			fmt.Printf("Error handling User-Agent: %s\n", err)
			return
		}
	} else if prefix := "/echo/"; strings.HasPrefix(requestTarget, prefix) {
		const acceptEncodingPrefix = "Accept-Encoding: "
		acceptEncoding, err := getHeaderValue(acceptEncodingPrefix, stringBuffer)

		var contentEncodingHeader string
		if err == nil && acceptEncoding == "gzip" {
			contentEncodingHeader = "Content-Encoding: gzip\r\n"
		}

		echoBody := strings.TrimPrefix(requestTarget, prefix)
		contentLength := len(echoBody)
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n%sContent-Length: %d\r\n\r\n%s", contentEncodingHeader, contentLength, echoBody)
	} else if prefix := "/files/"; strings.HasPrefix(requestTarget, prefix) {
		fileName := strings.TrimPrefix(requestTarget, prefix)
		filePath := directory + "/" + fileName

		response, err = handleFiles(requestMethod, stringBuffer, filePath)
		if err != nil {
			fmt.Printf("Error handling files: %s\n", err)
			return
		}
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("Error writing response: %s\n", err)
	}
}
