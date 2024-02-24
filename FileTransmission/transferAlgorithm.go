package filetransmission

import (
	helper "CloudDrive/Helper"
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

// File sizes
const (
	_  = iota
	KB = uint32(1 << (10 * iota)) // 1 KB = 1000 bytes
	MB = uint32(1 << (10 * iota)) // 1 MB = 1000000 bytes
	GB = uint32(1 << (10 * iota)) // 1 GB = 1000000000 bytes

	DrivePath = "D:\\CloudDrive"
)

// Based on our research for optimizing, returns the best chunk size (amount of bytes) for a given file size
func GetChunkSize(fileSize uint32) uint {

	switch {
	case fileSize < 10*MB:
		return uint(fileSize)
	case fileSize < 100*MB:
		return 256
	case fileSize < 400*MB:
		return 768
	case fileSize < 700*MB:
		return 1024
	case fileSize < 1*GB:
		return 2048
	default:
		return 0
	}
}

// Upload a file to the client
func SendFile(conn *net.Conn, size uint32, path string, name string) error {
	file, err := os.Open(DrivePath + path)

	if err != nil {
		return err
	}

	defer file.Close()

	// Recieve the chunksize for the uploaded file size
	chunkSize := GetChunkSize(size)
	chunk := make([]byte, chunkSize) // Makes a slice of bytes in size of the chunk size

	for { // Reads the file
		bytesRead, err := file.Read(chunk)
		if err == io.EOF || bytesRead == 0 {
			break
		}
		if err != nil {
			return err
		}
		err = helper.SendData(conn, chunk) // Sends the file data according to the chunk size
		if err != nil {
			return err
		}
	}
	return nil
}

// Reccive file from the client, read the data in chunks and then write in chunks into the file.
func ReceiveFile(conn net.Conn, filePath string, fileName string, fileSize int) error {
	fileBytes := make([]byte, fileSize) // Save the file content on a chunk bytes
	bytesRead := 0
	for bytesRead < fileSize { // First reading the file (to make sure the entire chunks can be read before writing to file)
		// Reading file
		read, err := conn.Read(fileBytes[bytesRead:])
		if err != nil {
			return fmt.Errorf("error reading the file.\nplease try again")
		}
		bytesRead += read
	}
	fullPath := filePath + "\\" + fileName
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0644) // Open file for writing
	if err != nil {
		return fmt.Errorf("error opening the file: %v", fileName)
	}
	// Create a buffered writer for efficient writes
	writer := bufio.NewWriter(file)

	//fileBytes := make([]byte, fileSize) // Save the file content on a chunk bytes
	bytesWritten := 0
	for bytesWritten < len(fileBytes) {
		n, err := writer.Write(fileBytes[bytesWritten:])
		if err != nil {
			return fmt.Errorf("error writing to file: %v", fileName)
		}
		bytesWritten += n
	}
	fmt.Println("Finished writing file")

	err = writer.Flush() // Flush any remaining data in the buffer to the file
	if err != nil {
		return fmt.Errorf("error flushing data to the file: %v", fileName)
	}

	return nil
}

/*

	for bytesRead < fileSize { // First reading the file (to make sure the entire chunks can be read before writing to file)
		read, err := conn.Read(fileBytes[bytesRead:])
		if err != nil {
			return fmt.Errorf("error reading the file.\nplease try again")
		}

		bytesRead += read
	}

	// Write the data to the file in chunks
	bytesWritten := 0
	for bytesWritten < len(fileBytes) {
		n, err := writer.Write(fileBytes[bytesWritten:])
		if err != nil {
			return fmt.Errorf("error writing to file: %v", fileName)
		}
		bytesWritten += n
	}

	err = writer.Flush() // Flush any remaining data in the buffer to the file
	if err != nil {
		return fmt.Errorf("error flushing data to the file: %v", fileName)
	}

	return nil*/
