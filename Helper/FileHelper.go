package helper

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DrivePath = "C:\\CloudDrive"
	RootDir   = "root"
)

func BuildUserFileSystem(userID uint32) error {
	rootPath := GetUserStorageRoot(userID)

	err := os.Mkdir(rootPath, os.ModePerm) // sets permissions for the directory
	if err != nil {
		return err
	}

	err = os.Mkdir(filepath.Join(rootPath, "Garbage"), os.ModePerm)
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
	return filepath.Join(DrivePath, strconv.FormatUint(uint64(userID), 10))
}
func GetUserStoragePath(userID uint32, clientPath string) string {
	return filepath.Join(GetUserStorageRoot(userID), clientPath[4:])
}

// need to fix
func GetVirtualStoragePath(userID uint32, storagePath string) string {
	return strings.Replace(storagePath, GetUserStorageRoot(userID), RootDir, -1)
}
