package FileSystem

import (
	helper "CloudDrive/Helper"
	"os"
	"path/filepath"
)

func (user *LoggedUser) ChangeDirectory(parameter string) error {

	switch parameter {
	case "\\":
		return user.setBackRoot()
	case "..":
		return user.setBackDirectory()
	default:
		if filepath.IsAbs(parameter) {
			return user.setAbsDir(parameter)
		}
		return user.setForwardDir(parameter)
	}

}

func (user *LoggedUser) CreateFile(fileName string) error {
	if !filepath.IsAbs(fileName) {
		fileName = filepath.Join(user.CurrentPath, fileName)
	}
	return createAbsFile(fileName)
}

func (user *LoggedUser) CreateDirectory(directoryName string) error {
	if !filepath.IsAbs(directoryName) {
		directoryName = filepath.Join(user.CurrentPath, directoryName)
	}
	return createAbsDir(directoryName)
}

func (user *LoggedUser) RenameFile(oldFileName string, newFileName string) error {
	if !filepath.IsAbs(oldFileName) {
		oldFileName = filepath.Join(user.CurrentPath, oldFileName)
	}
	if !filepath.IsAbs(newFileName) {
		newFileName = filepath.Join(user.CurrentPath, newFileName)
	}
	return renameAbsFile(oldFileName, newFileName)
}

// go back to the root, like cd/
func (user *LoggedUser) setBackRoot() error {
	return user.SetPath(helper.GetUserStorageRoot(user.UserID))
}

func (user *LoggedUser) setBackDirectory() error {
	newPath := filepath.Dir(user.CurrentPath)
	if newPath == "." {
		return &PathNotExistError{newPath}
	}
	return user.SetPath(newPath)
}

// forward to next given dir, like cd homework11
func (user *LoggedUser) setForwardDir(forwardDir string) error {
	err := isFileInDirectory(forwardDir, user.CurrentPath)
	if err != nil {
		return err
	}
	return user.SetPath(filepath.Join(user.CurrentPath, forwardDir))
}

func (user *LoggedUser) setAbsDir(absDir string) error {
	return user.SetPath(absDir)
}

// creat a file in the given directory
func createAbsFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func createAbsDir(absDir string) error {
	err := os.Mkdir(absDir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// rename a file name in the current directory, gets full file path and new filename
func RenameAbsFile(currentFilePath, newFileName string) error {
	newFileName = filepath.Join(filepath.Dir(currentFilePath) + newFileName)
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
