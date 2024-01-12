package Filesystem

import (
	"errors"
	"fmt"
	"os"
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
		return errors.New("Path is invalid")
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
	rootPath := DrivePath + rootName
	err := os.MkdirAll(rootPath, 0755) // 0755 sets permissions for the directory
	if err != nil {
		return err
	}
	drivePath := rootPath + "MyDrive\\"
	err = os.MkdirAll(drivePath, 0755)
	if err != nil {
		return err
	}

	garbagePath := rootPath + "Garbage\\"
	// sub-folder represent the folder that contins all the user data
	err = os.MkdirAll(garbagePath, 0755)
	if err != nil {
		return err
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
