package FileSystem

import (
	helper "CloudDrive/Helper"
	"errors"
	"io/fs"
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
			if errors.Is(err, fs.ErrNotExist) || errors.Is(err, fs.ErrExist) || errors.Is(err, fs.ErrPermission) || errors.Is(err, fs.ErrInvalid) { // If caught filesystem errors
				return nil, &InitializeError{}
			}
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

// Returns the user's current directory
func (user *LoggedUser) GetPath() string {
	return user.CurrentPath
}

// Returns the root total size
func (user *LoggedUser) GetRootSize() (uint32, error) {
	var totalSize uint32 = 0
	rootPath := helper.GetUserStorageRoot(user.UserID)

	// Traverses the root directory and calls the provided function for each file or directory encountered.
	err := filepath.Walk(rootPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() { // If the current item is a file, add its size to the total.
			totalSize += uint32(info.Size())
		}
		// Continue the traversal.

		return nil
	})
	if err != nil {
		return 0, &SizeCalculationError{}
	}
	return totalSize, nil

}

func isFolderInDirectory(path, pathOfDir string) error {
	userPath := path
	if !strings.Contains(path, helper.DrivePath) { // If path is not fully absolute
		userPath = filepath.Join(pathOfDir, path) // Convert path to fully absolute
	}
	_, err := os.Stat(userPath)

	if os.IsNotExist(err) {
		return &PathNotExistError{userPath}
	} else if err != nil {
		return err
	}

	return nil
}

// Valid the path that is given in the parameters of the client
func ValidPath(userID uint32, path string) error {
	if !strings.HasPrefix(path, helper.GetUserStorageRoot(userID)) {
		return &PremmisionError{path}
	}
	_, err := os.Stat(path)

	if err != nil { // If invalid path
		if os.IsNotExist(err) {
			return &PathNotExistError{path}
		}
		return err
	}

	return nil
}
