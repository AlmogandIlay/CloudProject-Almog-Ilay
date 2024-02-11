package FileSystem

import (
	helper "CloudDrive/Helper"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileManager: API interface for LoggedUser to interact with the file commands in the cloud drive.

// TODO: relate to the garbage in all the functions

// Changes the current directory for the user according to the parameter
func (user *LoggedUser) ChangeDirectory(parameter string) (string, error) {

	switch parameter {
	case "\\", "/":
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

// Creates a new file in the current directory
func (user *LoggedUser) CreateFile(fileName string) error {
	return user.fileOperation(fileName, createAbsFile)
}

// Creates a new folder in the current directory
func (user *LoggedUser) CreateFolder(folderName string) error {
	return user.fileOperation(folderName, createAbsDir)
}

// Remove a file in the current directory
func (user *LoggedUser) RemoveFile(fileName string) error {
	return user.fileOperation(fileName, removeAbsFile)
}

// Remove a folder in the current directory
func (user *LoggedUser) RemoveFolder(folderName string) error {
	return user.fileOperation(folderName, removeAbsFolder)
}

// Renames a file
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

// file operation that recieves all the file operations and send them to the function that is responsible for
func (user *LoggedUser) fileOperation(path string, operation func(string) error) error {
	if !filepath.IsAbs(user.CurrentPath) {
		path = filepath.Join(user.CurrentPath, path) // convert to an absolute path
	}
	return operation(path)
}

// go back to the CloudDrive user root: Root/
func (user *LoggedUser) setBackRoot() (string, error) {
	err := user.SetPath(helper.GetUserStorageRoot(user.UserID))
	return user.GetPath(), err
}

// Go back to the previous folder in the CloudDrive path user root
func (user *LoggedUser) setBackDirectory() (string, error) {
	newPath := filepath.Dir(user.CurrentPath)
	if newPath == "." {
		return "", &PathNotExistError{newPath}
	}
	err := user.SetPath(newPath)
	return user.GetPath(), err
}

// forward to next given dir, for example: cd homework11
func (user *LoggedUser) setForwardDir(forwardDir string) (string, error) {
	err := isFolderInDirectory(forwardDir, user.CurrentPath)
	if err != nil {
		return "", err
	}
	err = user.SetPath(filepath.Join(user.CurrentPath, forwardDir))
	return user.GetPath(), err
}

func (user *LoggedUser) setAbsDir(absDir string) (string, error) {
	err := user.SetPath(absDir)
	return user.GetPath(), err
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

// Create an absolute directory with the parameter folder name
func createAbsDir(absDir string) error {
	err := os.Mkdir(absDir, os.ModePerm) // Create a folder
	if err != nil {
		if os.IsExist(err) { // If folder exists
			return &FileExistError{absDir, filepath.Dir(absDir)}
		}
		return err
	}
	return nil
}

// Remove an absolute file with the parameter file name
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

// Remove an absolute directory with the parameter folder name
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

// Rename a file name in the current directory, gets full file path and new filename
func renameAbsFile(currentFilePath, newFileName string) error {
	newFileName = filepath.Join(filepath.Dir(currentFilePath), newFileName)
	// Use os.Rename to rename the file
	err := os.Rename(currentFilePath, newFileName)
	if err != nil {
		return err
	}
	return nil
}

func moveFile(currentAbsFilePath, newFileName string) error {
	return renameAbsFile(currentAbsFilePath, newFileName)
}

// Returns folder's content including its files and folders in a string variable. Return error if it fails
func getFolderContent(folderPath string) (string, error) {

	folder, err := os.Open(folderPath)
	if err != nil {
		return "", err
	}
	defer folder.Close()

	entries, err := folder.Readdir(-1) // Saves all entries in the directory

	if err != nil {
		return "", &ReadDirError{folderPath}
	}

	var builder strings.Builder
	var fileCounter, dirCounter uint

	for _, entry := range entries {
		// "dd/mm/yyyy HH:mm:ss <DIR|FILE> <filename>      at the end i need count file and folder"
		if entry.IsDir() { // If the entry is a directory
			dirCounter++
			WriteString(&builder, " <DIR> ", entry.Name(), "\n")
		} else { // Else if the entry is a file
			fileCounter++
			WriteString(&builder, " <FILE> ", entry.Name(), "\n")
		}
	}
	WriteString(&builder, fmt.Sprint((dirCounter)), "Dir(s), ", fmt.Sprint((fileCounter)), "File(s)")

	return builder.String(), nil
}

// Append entry to the builder.string object contains all the entries of the directory
func WriteString(builder *strings.Builder, fileParameter ...string) {
	for _, stringObject := range fileParameter {
		builder.WriteString(stringObject)
	}
}
