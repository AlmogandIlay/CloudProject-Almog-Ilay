package FileSystem

import (
	"CloudDrive/helper"
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
	return user.fileOperation(fileName, createAbsFile)
}

func (user *LoggedUser) CreateFolder(folderName string) error {
	return user.fileOperation(folderName, createAbsDir)
}

func (user *LoggedUser) RemoveFile(fileName string) error {
	return user.fileOperation(fileName, removeAbsFile)
}

func (user *LoggedUser) RemoveFolder(folderName string) error {
	return user.fileOperation(folderName, removeAbsFolder)
}

func (user *LoggedUser) RenameFile(filePath string, newFileName string) error {
	if !filepath.IsAbs(filePath) {
		err := isFileInDirectory(filePath, user.CurrentPath)
		if err != nil {
			return err
		}
		filePath = filepath.Join(user.CurrentPath, filePath) // if the file not abs -> file.* -> patn/file.*
	}
	return renameAbsFile(filePath, newFileName)
}

func (user *LoggedUser) fileOperation(path string, operation func(string) error) error {
	if !filepath.IsAbs(user.CurrentPath) {
		path = filepath.Join(user.CurrentPath, path)
	}
	return operation(path)
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

// create a file in the given directory
func createAbsFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		if os.IsExist(err) {
			return &FileExistError{filePath, filepath.Dir(filePath)}
		}
		return err
	}
	defer file.Close()
	return nil
}

func createAbsDir(absDir string) error {
	err := os.Mkdir(absDir, os.ModePerm)
	if err != nil {
		if os.IsExist(err) {
			return &FileExistError{absDir, filepath.Dir(absDir)}
		}
		return err
	}
	return nil
}

func removeAbsFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &PathNotExistError{filePath}
		}
		return err
	}
	return nil
}

func removeAbsFolder(filePath string) error {
	err := os.RemoveAll(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &PathNotExistError{filePath}
		}
		return err
	}
	return nil
}

// rename a file name in the current directory, gets full file path and new filename
func renameAbsFile(currentFilePath, newFileName string) error {
	newFileName = filepath.Join(filepath.Dir(currentFilePath), newFileName)
	// Use os.Rename to rename the file
	err := os.Rename(currentFilePath, newFileName)
	if err != nil {
		return err
	}
	return nil
}
