package FileRequestsManager

import (
	"client/ClientErrors"
	"client/Helper"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	empty    = 0
	kilobyte = 1_000_000
)

// Uploads file to the cloud server
func uploadFile(fileSize int64, chunksSize int, filename string, socket net.Conn) {
	file, err := os.Open(strings.Replace(filename, "'", "", Helper.RemoveAll)) // Open file
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	chunk := make([]byte, chunksSize) // Save buffer of chunks

	var validUpload = true
	var totalBytesRead int64
	var totalReadFlag int64 // Flag for client view to automatically update upload percentage and bar progress
	var precentage int64

	for {
		bytesRead, err := file.Read(chunk)
		if err == io.EOF { // If finish reading file succesfully
			break
		}
		if err != nil { // If error occurred while reading the file
			validUpload = false
			fmt.Println(err.Error())
			err = &ClientErrors.ReadFileInfoError{}
			fmt.Println(err.Error())
		}
		if bytesRead == empty { // If finish reading file succesfully
			break
		}

		_, err = socket.Write(chunk[:bytesRead]) // Sending chunk to server
		if err != nil {                          // If sending error occured
			validUpload = false
			err = &ClientErrors.SendDataError{}
			fmt.Println(err.Error())
		}

		// Add bytes read for the chunk to the flag percentage and to the total read bytes
		totalReadFlag += int64(bytesRead)
		totalBytesRead += int64(bytesRead)

		if totalReadFlag >= kilobyte { // For every 1 Kilobyte update the progess and perecntage bar
			totalReadFlag = 0
			precentage = (totalBytesRead * 100) / fileSize // Calculates total read bytes compared to the total file size in percentages
			printer := func(length int, char string) string {
				var bar string
				for i := 0; i < int(precentage/2); i++ { // bar is 2 times smaller than the actual percentages
					bar += "-"
				}
				return bar
			}
			fmt.Printf("\r%v    Upload Progress: %% - %s", precentage, printer(int(precentage), "-"))
		}
	}
	if validUpload {
		fmt.Printf("File %s has been uploaded successfully\n", filename)
	}
}
