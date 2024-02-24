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
	addr       = "192.168.50.220:12345"
	uploadAddr = "192.168.50.220:12346"
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

func handleUploadFile(conn net.Conn) {
	defer conn.Close()
}
func handleConnection(conn net.Conn, uploadListner net.Listener) {
	defer conn.Close()

	// Initialize setup
	printNewRemoteAddr(conn)
	var userHandler RequestHandlers.IRequestHandler = initializeRequestHandler() // Initialize handler interface for requests
	var loggedUser FileSystem.LoggedUser                                         // Logged User initialize
	closeConnection := false

	for !closeConnection {

		request_Info, err := Requests.ReciveRequestInfo(&conn) // Recive request info from client
		if err != nil {
			closeConnection = true
		}
		if request_Info.Type == Requests.UploadFileRequest {
			uploadConn, err := uploadListner.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}
			go handleUploadFile(uploadConn)

			response_info := userHandler.HandleRequest(request_Info, &loggedUser, &uploadConn) // Handle request processing

			err = RequestHandlers.SendResponseInfo(&uploadConn, response_info) // Send Response Info to client
			if err != nil {                                                    // If sending request info was unsucessful
				closeConnection = true
			}
			userHandler = RequestHandlers.UpdateRequestHandler(response_info) // Update Request Handler if needed

		} else {
			response_info := userHandler.HandleRequest(request_Info, &loggedUser, &conn) // Handle request processing

			err = RequestHandlers.SendResponseInfo(&conn, response_info) // Send Response Info to client
			if err != nil {                                              // If sending request info was unsucessful
				closeConnection = true
			}
			userHandler = RequestHandlers.UpdateRequestHandler(response_info) // Update Request Handler if needed

		}

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

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	uploadListener, err := net.Listen("tcp", uploadAddr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer uploadListener.Close()

	fmt.Printf("Server is listening on %s...\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, uploadListener)
	}

}
