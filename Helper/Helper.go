package helper

import (
	"fmt"
	"net"
	"server/RequestHandlers"
)

type IRequestHandler interface {
	RequestHandlers.IRequestHandler
}

// bufferSize is usually 1024
func ReciveData(conn *net.Conn, bufferSize int) []byte {
	buffer := make([]byte, bufferSize)
	_, err := (*conn).Read(buffer)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return nil
	}
	return buffer
}

func SendData(conn *net.Conn, message []byte) {
	_, err := (*conn).Write(message)
	if err != nil {
		fmt.Println("Error sending response:", err)
	}
}
