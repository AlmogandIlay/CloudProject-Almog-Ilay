package main

import (
	"CloudDrive/FileSystem"
	"CloudDrive/Server/RequestHandlers"
	"CloudDrive/Server/RequestHandlers/Requests"
	"fmt"
	"log"
	"net"
)

const (
	addr = "192.168.50.220:12345"
)

//Prints the Remote IP:Port's client in the CLI server program.

func printNewRemoteAddr(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), "is connected")
}

func printDisconnectedRemoteAddr(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), "is disconnected")
}

//Initializes RequestHandler variable

func initializeRequestHandler() RequestHandlers.AuthenticationRequestHandler {
	return RequestHandlers.AuthenticationRequestHandler{}
}

//Handles new client connection

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Initialize setup
	printNewRemoteAddr(conn)
	var userHandler RequestHandlers.IRequestHandler = initializeRequestHandler()
	var loggedUser FileSystem.LoggedUser
	closeConnection := false

	for !closeConnection {

		request_Info, err := Requests.ReciveRequestInfo(&conn)
		if err != nil {
			closeConnection = true
		}
		response_info := userHandler.HandleRequest(request_Info, &loggedUser)

		err = RequestHandlers.SendResponseInfo(&conn, response_info)
		if err != nil { // If sending request info was unsucessful
			closeConnection = true
		}
		userHandler = RequestHandlers.ChangeRequestHandler(response_info)

	}

	printDisconnectedRemoteAddr(conn)
}

func main() {
	_, err := RequestHandlers.InitializeIdentifyManagerFactory()
	if err != nil {
		log.Fatal("There has been an error when attempting to initialize Factory.\nError Data:", err.Error())
		return
	}

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
