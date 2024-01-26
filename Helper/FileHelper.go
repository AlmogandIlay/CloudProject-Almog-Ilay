package helper

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	DrivePath = "C:\\CloudDrive\\"
	RootDir   = "root\\"
)

func BuildUserFileSystem(userID uint32) error {
	rootPath := fmt.Sprintf("%s%s%s", DrivePath, fmt.Sprint(userID), "\\")

	err := os.MkdirAll(rootPath, 0755) // 0755 sets permissions for the directory
	if err != nil {
		return err
	}

	err = os.MkdirAll(rootPath+"Garbage\\", 0755)
	if err != nil {
		return err
	}
	return nil
}

/*
from the user observation he see and write the encapsulate path(root/filesys/etc...) but the server validate the literal path(C:/CloudDrive)
root/file1/file2 -> C:/CloudDrive/id/file1/file2
*/

func GetUserStorageRoot(userID uint32) string {

	return fmt.Sprint(DrivePath, strconv.FormatUint(uint64(userID), 10))
}
func GetStoragePath(userID uint32, clientPath string) string {
	return GetUserStorageRoot(userID) + clientPath[4:]
}

func GetVirtualStoragePath(userID uint32, storagePath string) string {
	return strings.Replace(storagePath, GetUserStorageRoot(userID), RootDir, -1)
}
