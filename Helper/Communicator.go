package helper

import (
	"net"
	"strings"
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
