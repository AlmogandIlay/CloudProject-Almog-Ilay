package FileSystem

import (
	helper "CloudDrive/Helper"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// go back to the root, like cd/
func (user *LoggedUser) SetBackRoot() error {
	return user.SetPath(helper.GetUserStorageRoot(user.UserID))
}

func (user *LoggedUser) SetBackwardDir() error {
	newPath := filepath.Dir(user.CurrentPath)
	if newPath == "." {
		return &PathNotExistError{newPath}
	}
	return user.SetPath(newPath)
}

// forward to next dir, like cd homework11
func (user *LoggedUser) SetForwardDir(forwardDir string) error {
	err := isFileInDirectory(forwardDir, user.CurrentPath)
	if err != nil {
		return err
	}
	fmt.Println("filefile")
	return user.SetPath(user.CurrentPath + "\\" + forwardDir)
}

// not finish, need to check for forward directory
func (user *LoggedUser) SetAbsDir(absDir string) error {

	if strings.Contains(user.CurrentPath, "root\\") {
		return errors.New("1")
	}

	tempPath := filepath.Dir(user.CurrentPath) // cd..
	lastDir := filepath.Base(user.CurrentPath) // return the current dir name

	for {
		err := isFileInDirectory(lastDir, tempPath)
		if err != nil {
			return err
		}

		if tempPath == "root\\" {
			return user.SetPath(absDir)
		}

		tempPath = filepath.Dir(tempPath)
		lastDir = filepath.Base(lastDir)
	}
}

/*
func (user *LoggedUser) GetStorage() (uint, error) {
	userDrivePath := helper.DrivePath + "\\" + user.UserID


}*/

// rename a file name in the current directory

func RenameAbsFile(currentFilePath, newFileName string) error {

	newFileName = filepath.Dir(currentFilePath) + "\\" + newFileName
	// Use os.Rename to rename the file
	err := os.Rename(currentFilePath, newFileName)
	if err != nil {
		return err
	}
	return nil
}

// on defualt: currentUserPath = loggedUser.CurrentPath
func RenameRelativeFile(currentFileName, newFileName, currentUserPath string) error {
	err := isFileInDirectory(currentFileName, currentUserPath)

	if err != nil {
		return err
	}
	currentFileName = currentUserPath + "\\" + currentFileName
	return RenameAbsFile(currentFileName, newFileName)

}

func CreateAbsFile(fileName string) error {
	return nil
}
