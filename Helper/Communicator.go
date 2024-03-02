package Helper

import (
	"client/ClientErrors"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type AmountOfPaths int

const (
	DefaultBufferSize                 = 1024
	RemoveAll                         = -1
	FirstNameParameter                = 0
	SecondNameParameter               = 2
	OneClosedPath       AmountOfPaths = 2 // Specifcy that looking for one filename argument that is closed with ''
	TwoCloudPaths       AmountOfPaths = 4 // Specifcy that looking for two filename arguments that is closed with '' ''
	secondPathIndex                   = 3
	chunksIndex                       = 2
	SkipEnclose                       = 1

	uploadAddr = "clouddriveserver.duckdns.org:12346"
)

// bufferSize is usually 1024

/*
Recive data from client socket with the given buffer size.

Returns the received bytes.
*/
func ReciveData(conn *net.Conn, bufferSize int) (dataBytes []byte, errr error) {
	buffer := make([]byte, bufferSize)

	err := (*conn).SetReadDeadline(time.Now().Add(10 * time.Second)) // Set timeout for packet to recieve
	if err != nil {
		return nil, fmt.Errorf("it took too long time to get a respone back from the server")
	}
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

// Find if the parameters are closed with ' in the amount of the given closedCount (for example, given TwoClosedPath: 'filename1' 'filename2')
func IsQuoted(command_arguments []string, closedCount AmountOfPaths) bool {
	arguments := strings.Join(command_arguments, " ")
	var counter AmountOfPaths = 0
	for _, char := range arguments {
		if char == '\'' {
			counter++
		}
	}
	return counter == closedCount
}

// Find the path of a file if it's closed with ' with the given amount of closed paths and specific filename index.
// Returns the string path with enclouse '
func FindPath(command_arguments []string, index int, closedCount AmountOfPaths) string {
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
	if index == SecondNameParameter { // If the second path has been chosen
		parts := strings.Split(name, "'")
		// If the client provided 2 paths (two single quotes), save the second part
		if len(parts) >= int(TwoCloudPaths) {
			name = parts[secondPathIndex]
		} else {
			name = ""
		}
		//name = "'" + strings.Join((strings.Split(name, "'")[1])) + "'" // save the second path only. (int[closedCount]-1) to access the last index of the closedCount path
	}
	return name
}

// For cases when the first path is quoted but the second doesn't, return the second Non-quoted second path
func ReturnNonQuotedSecondPath(command_arguments []string) string {
	arguments := strings.Join(command_arguments, " ") // Convert the comman arguments to string
	lastIndex := strings.LastIndex(arguments, "'")    // Find the last time the ' index encloused has been used in the string and append 1 to avoid it completely
	if lastIndex+1 <= len(arguments) {                // If Second path hasn't been specified
		return ""
	}
	return arguments[lastIndex+1:]
}

// Creates a private socket connection between the server for file transmission
func CreatePrivateSocket() (*net.Conn, error) {
	sock, err := net.Dial("tcp", uploadAddr)
	if err != nil {
		fmt.Println(err.Error())
		return nil, &ClientErrors.ServerConnectionError{Err: err}
	}
	return &sock, nil
}

// Convert respone to chunks size
func ConvertResponeToChunks(respone string) (int, error) {
	return strconv.Atoi(strings.Split(respone, ":")[chunksIndex])
}
