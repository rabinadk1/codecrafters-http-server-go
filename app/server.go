package main

import (
	"fmt"
	"log"
	"net"
)

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
			log.Fatalln("Error accepting connection: ", err.Error())
		}

		go func(conn net.Conn) {
			defer conn.Close()

			fmt.Printf("Accepted connection for %s\n", conn.RemoteAddr())

			_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
			if err != nil {
				log.Fatalf("Error writing response: %s\n", err)
			}
		}(conn)
	}
}
