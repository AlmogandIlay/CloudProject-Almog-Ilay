package helper

import (
	"net"
	"strings"
	"time"
)

const (
	DefaultBufferSize = 1024
)

// bufferSize is usually 1024

/*
Recive data from client socket with the given buffer size and with the timeout flag option.
Returns the received bytes or error.
*/
func ReciveData(conn *net.Conn, bufferSize int, timeoutFlag bool) ([]byte, error) {
	buffer := make([]byte, bufferSize)
	if timeoutFlag { // If timeout flag is on
		err := (*conn).SetReadDeadline(time.Now().Add(10 * time.Second)) // Set timeout for packet to recieve (10 seconds)
		if err != nil {
			return nil, err
		}
	}
	bytesRead, err := (*conn).Read(buffer)
	if err != nil {
		return nil, err
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
		return err
	}
	return nil
}

func ConvertRawJsonToData(rawJson string) string {
	data := strings.ReplaceAll(rawJson, "\"", "")
	data = strings.ReplaceAll(data, "}", "")
	return data[6:]
}
