package main

import (
	helper "CloudDrive/Helper"
	"CloudDrive/Server/RequestHandlers"
	"fmt"
	"log"
	"net"
)

/*
Prints the Remote IP:Port's client in the CLI server program.
*/
func printNewRemoteAddr(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), "is connected")
}

func printDisconnectedRemoteAddr(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), "is disconnected")
}

/*
Initializes RequestHandler variable
*/
func initializeRequestHandler() *RequestHandlers.AuthenticationRequestHandler {
	return &RequestHandlers.AuthenticationRequestHandler{}
}

/*
Handles new client connection
*/
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Initialize setup
	printNewRemoteAddr(conn)
	userHandler := initializeRequestHandler()
	closeConnection := false

	for !closeConnection {

		requestInfo, err := helper.ReciveRequestInfo(&conn)
		if err != nil {
			closeConnection = true
		}
		err = helper.SendResponseInfo(&conn, userHandler.HandleRequest(requestInfo))

		if err != nil { // If sending request info was unsucessful
			closeConnection = true
		}
	}

	printDisconnectedRemoteAddr(conn)
}

func main() {
	_, err := RequestHandlers.InitializeFactory()
	if err != nil {
		log.Fatal("There has been an error when attempting to initialize Factory.\nError Data:", err.Error())
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
