package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	directory := "."
	if len(os.Args) > 1 && os.Args[1] == "--directory" {
		if len(os.Args) < 3 {
			log.Fatalln("Usage: ./http_server --directory <directory>")
		}
		directory = strings.TrimRight(os.Args[2], "/")
	}

	const port = 4221
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("Failed to bind to port %d: %s\n", port, err)
	}
	defer l.Close()

	fmt.Printf("Listening on port %d\n", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}

		go handleConnection(conn, directory)

	}
}
