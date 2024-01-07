package helper

import (
	"CloudDrive/Server/RequestHandlers"
	"CloudDrive/Server/RequestHandlers/Requests"
	"encoding/json"
	"fmt"
	"net"
)

const (
	defaultBufferSize = 1024
)

// bufferSize is usually 1024

/*
Recive data from client socket with the given buffer size.

Returns the received bytes.
*/
func ReciveData(conn *net.Conn, bufferSize int) []byte {
	buffer := make([]byte, bufferSize)
	bytesRead, err := (*conn).Read(buffer)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return nil
	}

	data := make([]byte, bytesRead)
	copy(data, buffer[:bytesRead])
	return data
}

// Recives data from client socket and returns RequestInfo
func ReciveRequestInfo(conn *net.Conn) (Requests.RequestInfo, error) {
	data := ReciveData(conn, defaultBufferSize)

	var requestInfo Requests.RequestInfo

	err := json.Unmarshal(data, &requestInfo)
	if err != nil {
		return Requests.RequestInfo{Requests.ErrorRequest, []byte("Error: One or more of the fields are being used wrong.")}, err
	}

	return requestInfo, nil

}

/*
Send data to the client socket.

Returns the received bytes.
*/
func SendData(conn *net.Conn, message []byte) error {
	_, err := (*conn).Write(message)
	if err != nil {
		return err
	}
	return nil
}

func SendResponseInfo(conn *net.Conn, responseInfo RequestHandlers.ResponeInfo) error {
	message, err := json.Marshal(responseInfo)
	if err != nil {
		return err
	}

	return SendData(conn, message)
}
