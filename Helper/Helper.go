package helper

import (
	"CloudDrive/Server/RequestHandlers"
	"CloudDrive/Server/RequestHandlers/Requests"
	"encoding/json"
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

// Recives data from client socket and returns RequestInfo
func ReciveRequestInfo(conn *net.Conn) (Requests.RequestInfo, error) {
	data, err := ReciveData(conn, defaultBufferSize)

	if err != nil {
		return Requests.RequestInfo{}, err
	}

	var requestInfo Requests.RequestInfo

	err = json.Unmarshal(data, &requestInfo)
	if err != nil {
		return Requests.RequestInfo{
			Type:        Requests.ErrorRequest,
			RequestData: []byte("Error: One or more of the fields are being used wrong."),
		}, nil

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
