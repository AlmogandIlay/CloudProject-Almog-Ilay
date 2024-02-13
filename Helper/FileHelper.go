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
func GetServerStoragePath(userID uint32, clientPath string) string {
	findAbsolute := strings.Index(clientPath, "\\")
	if findAbsolute == -1 {
		return clientPath
	}
	serverPath := DrivePath + "\\" + strconv.FormatUint(uint64(userID), 10) + clientPath[findAbsolute:]
	return serverPath
}

// Converts server-side path to client-side path
func GetVirtualStoragePath(storagePath string) string {
	var clientPath string

	cloudDriveIndex := strings.Index(storagePath, CloudDrive)

	// Extract the part after "CloudDrive", including separators
	pathAfterCloudDrive := storagePath[cloudDriveIndex:]

	// Split by '\'
	pathParts := strings.SplitN(pathAfterCloudDrive, "\\", 3)

	// discard the assumed user ID
	if len(pathParts) > 2 { // If the given path is expandable more than the root
		relativePath := pathParts[2]
		clientPath = fmt.Sprintf(RootDir+"%s", relativePath)
	} else { // If the given path is directly the root
		clientPath = RootDir
	}
	return clientPath

}
