package helper

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	CloudDrive = "CloudDrive"
	DrivePath  = "D:\\CloudDrive"
	RootDir    = "Root:\\"
)

/*
from the user observation he see and write the encapsulate path(root/filesys/etc...) but the server validate the literal path(C:/CloudDrive)
root/file1/file2 -> C:/CloudDrive/id/file1/file2
*/

func GetUserStorageRoot(userID uint32) string {
	return filepath.Join(DrivePath, strconv.FormatUint(uint64(userID), 10))
}
func GetUserStoragePath(userID uint32, clientPath string) string {
	return filepath.Join(GetUserStorageRoot(userID), clientPath[4:])
}

// Converts server-side path to client-side path
func GetVirtualStoragePath(storagePath string) string {
	cloudDriveIndex := strings.Index(storagePath, CloudDrive)

	// Extract the part after "CloudDrive", including separators
	pathAfterCloudDrive := storagePath[cloudDriveIndex:]

	// Split by '\'
	pathParts := strings.SplitN(pathAfterCloudDrive, "\\", 3)

	// discard the assumed user ID
	relativePath := pathParts[2]
	clientPath := fmt.Sprintf("Root\\%s", relativePath)
	return clientPath

}
