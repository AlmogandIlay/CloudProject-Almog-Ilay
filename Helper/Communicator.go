package helper

import (
	"fmt"
	"net"
)

const (
	DefaultBufferSize = 1024
)

// bufferSize is usually 1024

/*
Recive data from client socket with the given buffer size.

Returns the received bytes.
*/
func ReciveData(conn *net.Conn, bufferSize int) ([]byte, error) {
	buffer := make([]byte, bufferSize)
	bytesRead, err := (*conn).Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("error when reciving a response from the server.\nPlease send this info to the developers:\n%s", err)
	}

	data := make([]byte, bytesRead)
	copy(data, buffer[:bytesRead])
	return data, nil
}

/*
Send data to the client socket.

Returns the received bytes.
*/
func SendData(conn *net.Conn, message []byte) error {

	_, err := (*conn).Write(message)
	if err != nil {
		return fmt.Errorf("error when attempting to send the request to the server.\nPlease send this info to the developers:\n%s", err)
	}
	return nil
}
