package Filesystem

import (
	helper "CloudDrive/Helper"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type LoggedUser struct { // Might want to replace to a better name
	UserID      uint32
	CurrentPath string
}

func newLoggedUser(id uint32, path string) (*LoggedUser, error) {
	var user LoggedUser

	err := user.setPath(path)
	if err != nil {
		return nil, err
	}
	return &LoggedUser{id, path}, nil
}

func (user *LoggedUser) setPath(path string) error {
	if err := Valid(path); err != nil {
		return err
	}
	user.CurrentPath = path
	return nil
}

func Valid(path string) error {
	dir, err := os.Open(string(path))

	defer dir.Close()

	if os.IsNotExist(err) {
		return &PathNotExistError{path}
	}
	if err != nil { // If invalid path
		return err
	}

	return nil
}

func BuildUserFileSystem(userID uint32) error {

	rootName := fmt.Sprint(userID) + "\\"
	rootPath := helper.DrivePath + rootName
	err := os.MkdirAll(rootPath, 0755) // 0755 sets permissions for the directory
	if err != nil {
		return err
	}

	err = os.MkdirAll(rootPath+"MyDrive\\", 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(rootPath+"Garbage\\", 0755)
	if err != nil {
		return err
	}
	return nil
}

func (user *LoggedUser) setBackRoot() error {
	err := user.setPath("root\\")
	return err
}
func (user *LoggedUser) setForwardDir(forwardDir string) error {
	err := isFileInDirectory(forwardDir, user.CurrentPath)
	if err != nil {
		return err
	}
	return user.setPath(user.CurrentPath + forwardDir + "\\")
}

func (user *LoggedUser) setAbsDir(absDir string) error {

	if strings.Contains(user.CurrentPath, "root\\") {
		return errors.New("1")
	}

	tempPath := filepath.Dir(user.CurrentPath)
	lastDir := filepath.Base(user.CurrentPath)

	for {
		err := isFileInDirectory(lastDir, tempPath)
		if err != nil {
			return err
		}

		if tempPath == "root\\" {
			return user.setPath(absDir)
		}

		tempPath = filepath.Dir(tempPath)
		lastDir = filepath.Base(lastDir)
	}

}

// from the user observation he see and write the encapsulate path(root/filesys/etc...) but the server validate the literal path(C:/CloudDrive)
func getCloudPath(userID uint32, clientPath string) string {
	return DrivePath + fmt.Sprint(userID) + clientPath[4:]
}

func isFileInDirectory(fileName, pathOfDir string) error {
	dir, err := os.Open(pathOfDir)
	defer dir.Close()

	if err != nil {
		return err
	}

	files, err := dir.ReadDir(-1) // Saves all the path files in a slice
	if err != nil {
		return err
	}

	// check for file with same name. note: files are case-senstive
	for _, file := range files {
		if strings.ToLower(file.Name()) == strings.ToLower(fileName) {
			return err
		}
	}

	return nil
}

type FileOperation interface {
	Download() error //
	Upload() error   //
	Rename() error
	Create() error
	Delete() error
	TransferToGarbege() error
}
type FileSystemNode interface {
	FileOperation
	IsFile() bool
	IsFolder() bool
}
