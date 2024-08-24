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
