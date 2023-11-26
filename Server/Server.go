package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func handleConnection(conn net.Conn) {
	const SizeDigits = 13
	fmt.Println(conn.RemoteAddr(), "is connected")
	buf := make([]byte, SizeDigits)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	var size int
	var convert string
	for i := 0; i < SizeDigits; i++ {
		if buf[i] == 0 {
			break
		}
		convert += string(buf[i])
	}

	size, _ = strconv.Atoi(convert)
	buf = make([]byte, size)
	start := time.Now()
	_, err = conn.Read(buf)

	elapsed := time.Since(start)

	if err != nil {
		fmt.Println(err)
	}
	println("File content is:")
	println(string(buf))
	fmt.Println("\nSize of the given file is", size)
	fmt.Println("Time:", elapsed)
}

func main() {
	addr := "192.168.50.191:12345"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server is listening on %s...\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
