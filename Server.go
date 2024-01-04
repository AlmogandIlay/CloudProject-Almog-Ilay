package main

import (
	"CloudDrive/Server/RequestHandlers"
	"CloudDrive/Server/RequestHandlers/Requests"
	"fmt"
	"log"
	"net"
)

func printRemoteAddr(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), "is connected")
}

func handleConnection(conn net.Conn) {
	printRemoteAddr(conn)

	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	data := make([]byte, bytesRead)
	copy(data, buffer[:bytesRead])

	info := Requests.RequestInfo{0, data}
	handler := new(RequestHandlers.LoginRequestHandler)
	handler.HandleLogin(info)
}

func main() {
	_, err := RequestHandlers.InitializeFactory()
	if err != nil {
		log.Fatal("There has been an error when attempting to initializeFactory.\nError Data:", err.Error())
		return
	}

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
