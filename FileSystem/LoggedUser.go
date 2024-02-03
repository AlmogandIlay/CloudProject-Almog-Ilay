package FileSystem

import (
	helper "CloudDrive/Helper"
	"os"
	"path/filepath"
	"strings"
)

type LoggedUser struct {
	UserID      uint32
	CurrentPath string // server prespective path
}

func NewLoggedUser(id uint32) (*LoggedUser, error) {

	user := &LoggedUser{UserID: id}
	err := BuildUserFileSystem(id)
	if err != nil {
		return nil, err
	}
	err = user.SetPath(helper.GetUserStorageRoot(id))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func BuildUserFileSystem(userID uint32) error {
	rootPath := helper.GetUserStorageRoot(userID)

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

func (user *LoggedUser) SetPath(path string) error {
	err := ValidPath(user.UserID, path)
	if err != nil {
		return err
	}
	user.CurrentPath = path
	return nil
}

func ValidPath(userID uint32, path string) error {
	if !strings.HasPrefix(path, helper.GetUserStorageRoot(userID)) {
		return &PathNotExistError{path}
	}
	dir, err := os.Open(path)

	if err != nil { // If invalid path
		if os.IsNotExist(err) {
			return &PathNotExistError{path}
		}
		return err
	}

	defer dir.Close()
	return nil
}
