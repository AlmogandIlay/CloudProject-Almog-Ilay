package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), "is connected")
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	for i := 0; i < 1024; i++ {
		if buf[i] != 0 {
			fmt.Print(string(buf[i]))
		} else {
			break
		}
	}
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