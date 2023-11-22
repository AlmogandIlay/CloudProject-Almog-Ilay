package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func handleConnection(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), "is connected")

	const SizeDigits = 13
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
	fmt.Println("Size of the given file is", size)

	file := make([]byte, size)
	bytesRead := 0
	start := time.Now()
	for bytesRead < size {
		read, err := conn.Read(file[bytesRead:])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			fmt.Print(read)
		}
		bytesRead += read
	}
	fmt.Println("File has fully received!\nContent:")
	fmt.Println(string(file[:]))
	elapsed := time.Since(start)
	fmt.Printf("\n\nBinomial took %s\n\n", elapsed)

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
