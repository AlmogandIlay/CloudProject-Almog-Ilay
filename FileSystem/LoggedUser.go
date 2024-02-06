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

// Initializes a new logged user
func NewLoggedUser(id uint32) (*LoggedUser, error) {
	var err error
	user := &LoggedUser{UserID: id} // Creates a new user instance
	if !IsFileSystemExist(id) {     // Checks if the File System was not made before
		err = BuildUserFileSystem(id)
		if err != nil {
			return nil, err
		}
	}
	err = user.SetPath(helper.GetUserStorageRoot(id)) // Sets user's path
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Checks if the File System was created before
func IsFileSystemExist(userID uint32) bool {
	rootPath := helper.GetUserStorageRoot(userID)
	_, err := os.OpenFile(rootPath, os.O_RDONLY, os.ModeDir) // Try to open the directory
	return err == nil                                        // Returns if path exists
}

// Builds a File System in the cloud for new clients
func BuildUserFileSystem(userID uint32) error {
	rootPath := helper.GetUserStorageRoot(userID)

	err := os.Mkdir(rootPath, os.ModePerm) // sets permissions for the directory
	if err != nil {
		return err
	}

	err = os.Mkdir(filepath.Join(rootPath, "Garbage"), os.ModePerm) // Creates the Garbage location
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

// Valid the path that is given in the parameters of the client
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
