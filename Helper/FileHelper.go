package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	CloudDrive = "CloudDrive"
	DrivePath  = "D:\\CloudDrive"
	RootDir    = "Root:\\"
	noAbsolute = -1
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
	if findAbsolute == noAbsolute || !strings.HasPrefix(clientPath, RootDir) { // If the path is not absolute or not starts with Root:\\
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

// convert to an absolute-server side path
func ConvertToAbsolute(fullpath, filePath string) string {
	return filepath.Join(fullpath, filePath)
}

// IsPathSeparator reports whether c is a directory separator character.
func isPathSeparator(c uint8) bool {
	return c == '\\'
}

// Base returns the last element of path.
// Trailing path separators are removed before extracting the last element.
// If the path is empty, Base returns ".".
// If the path consists entirely of separators, Base returns a single separator.
func Base(path string) string {
	if path == "" {
		return "."
	}
	// Strip trailing slashes.
	for len(path) > 0 && os.IsPathSeparator(path[len(path)-1]) {
		path = path[0 : len(path)-1]
	}
	// Throw away volume name
	path = path[len(filepath.VolumeName(path)):]
	// Find the last element
	i := len(path) - 1
	for i >= 0 && !isPathSeparator(path[i]) {
		i--
	}
	if i >= 0 {
		path = path[i+1:]
	}
	// If empty now, it had only slashes.
	if path == "" {
		return string(filepath.Separator)
	}
	return path
}
