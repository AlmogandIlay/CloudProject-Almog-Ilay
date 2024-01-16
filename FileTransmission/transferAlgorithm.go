package filetransmission

import (
	"CloudDrive/Filesystem"
	"CloudDrive/helper"
	"io"
	"net"
	"os"
)

// add comment
const (
	_  = iota
	KB = uint32(1 << (10 * iota)) // 1 << (10 * 1)
	MB = uint32(1 << (10 * iota))
	GB = uint32(1 << (10 * iota))

	DrivePath = "C:\\CloudDrive"
)

// nununu Add comments here
func GetChunkSize(fileSize uint32) uint {

	switch {
	case fileSize < 10*MB:
		return uint(fileSize)
	case fileSize < 100*MB:
		return 0
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

// nununu add comments here
func SendFile(conn *net.Conn, uploadedFile *Filesystem.File) error {
	file, err := os.Open(DrivePath + uploadedFile.Path)

	if err != nil {
		return err
	}

	defer file.Close()

	// add comment
	chunkSize := GetChunkSize(uploadedFile.Size)
	chunk := make([]byte, chunkSize)

	for { // add comment
		bytesRead, err := file.Read(chunk)
		if err == io.EOF || bytesRead == 0 {
			break
		}
		if err != nil {
			return err
		}
		err = helper.SendData(conn, chunk)
		if err != nil {
			return err
		}
	}
	return nil
}
