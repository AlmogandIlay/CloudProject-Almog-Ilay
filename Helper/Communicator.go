package Helper

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

const (
	DefaultBufferSize   = 1024
	FirstNameParameter  = 0
	SecondNameParameter = 2
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

func ConvertStringToBytes(data string) ([]byte, error) {
	// Create a simple struct to hold data
	jsonData := struct {
		Data string `json:"Data"`
	}{Data: data}

	bytes, err := json.Marshal(jsonData)
	if err != nil {
		// Handle the error appropriately
		return nil, fmt.Errorf("error encoding requested path to be sent to the server")
	}
	return bytes, nil
}

// Find if the parameters are closed with ' (for example: 'filename1' 'filename2')
func IsQuoted(command_arguments []string) bool {
	arguments := strings.Join(command_arguments, " ")
	counter := 0
	for _, char := range arguments {
		if char == '\'' {
			counter++
		}
	}
	return counter == 4
}

// Find the path of a file if it's closed with '
func FindPath(command_arguments []string, index int) string {
	arguments := strings.Join(command_arguments, " ") // Convert to string
	var name string
	var counter int // counting the amount of '
	for _, char := range arguments {
		if char == '\'' { // If close ' char found
			counter++
		}
		name += string(char)    // append char to string
		if counter >= index+2 { // If 2/4 ' close chars have been found
			break // Stop searching
		}
	}
	if counter == 4 { // If the second path has been specifed
		name = "'" + strings.Split(name, "'")[3] + "'" // save the second path only
	}
	return name
}
