package main

import (
	"CloudDrive/FileSystem"
	"CloudDrive/Server/RequestHandlers"
	"CloudDrive/Server/RequestHandlers/Requests"
	"fmt"
	"log"
	"net"
)

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

var serverAddr string

var transmissionAddr string

func init() {
	serverAddr = getOutboundIP() + ":12345"
	transmissionAddr = getOutboundIP() + ":12346"
	fmt.Println(serverAddr)
}

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

func handleConnection(conn net.Conn, fileTransferListener net.Listener) {
	defer conn.Close()

	// Initialize setup
	printNewRemoteAddr(conn)
	var userHandler RequestHandlers.IRequestHandler = initializeRequestHandler() // Initialize handler interface for requests
	var loggedUser FileSystem.LoggedUser                                         // Logged User initialize
	closeConnection := false

	for !closeConnection {

		request_Info, err := Requests.ReciveRequestInfo(&conn, false) // Recive request info from client with no timeout flag
		if err != nil {
			closeConnection = true
		}
		response_info := userHandler.HandleRequest(request_Info, &loggedUser, &fileTransferListener) // Handle request processing

		err = RequestHandlers.SendResponseInfo(&conn, response_info) // Send Response Info to client
		if err != nil {                                              // If sending request info was unsucessful
			closeConnection = true
		}
		userHandler = RequestHandlers.UpdateRequestHandler(response_info) // Update Request Handler if needed
	}

	err := RequestHandlers.RemoveOnlineUser(loggedUser) // Remove the current user from the online users array
	if err != nil {
		return
	}
	printDisconnectedRemoteAddr(conn)
}

func main() {
	_, err := RequestHandlers.InitializeAuthenticationManagerFactory()
	if err != nil {
		log.Fatal("There has been an error when attempting to initialize Factory.\nError Data:", err.Error())
		return
	}

	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fileTransferListener, err := net.Listen("tcp", transmissionAddr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer fileTransferListener.Close()

	fmt.Printf("Server is listening on %s...\n", serverAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, fileTransferListener)
	}

}
