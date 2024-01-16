package helper

import "fmt"

const (
	DrivePath = "C:\\CloudDrive\\"
)

// from the user observation he see and write the encapsulate path(root/filesys/etc...) but the server validate the literal path(C:/CloudDrive)
func getCloudPath(userID uint32, clientPath string) string {
	return DrivePath + fmt.Sprint(userID) + clientPath[4:]
}
