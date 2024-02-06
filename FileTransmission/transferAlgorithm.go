package filetransmission

import (
	"CloudDrive/FileSystem"
	helper "CloudDrive/Helper"
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
	case fileSize < 900*MB:
		return 2048
	default:
		return 0
	}
}

// Upload a file to the client
func SendFile(conn *net.Conn, uploadedFile *FileSystem.File) error {
	file, err := os.Open(DrivePath + uploadedFile.Path)

	if err != nil {
		return err
	}

	defer file.Close()

	// Recieve the chunksize for the uploaded file size
	chunkSize := GetChunkSize(uploadedFile.Size)
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
