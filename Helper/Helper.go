package helper

import (
	"fmt"
	"net"
)

// bufferSize is usually 1024
func ReciveData(conn net.Conn, bufferSize int) []byte {
	buffer := make([]byte, bufferSize)
	_, err := (conn).Read(buffer)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return nil
	}
	return buffer
}

func SendData(conn net.Conn, message string) {
	_, err := (conn).Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending response:", err)
	}
}
