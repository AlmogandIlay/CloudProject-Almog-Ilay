package FileSystem

import (
	helper "CloudDrive/Helper"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileManager: API interface for LoggedUser to interact with the file commands in the cloud drive.

// TODO: relacte to the garbage in all the functions

// Changes the current directory for the user according to the parameter
func (user *LoggedUser) ChangeDirectory(parameter string) (string, error) {
	var err error
	var path string
	serverPath := helper.GetServerStoragePath(user.UserID, parameter)

	if serverPath != ".." && serverPath != "\\" && serverPath != "/" { // If the path is not going back or forward
		err = validFileName(helper.Base(serverPath), user.GetPath()) // Valid for files and folders are equals. calling the validFileName
		if err != nil {
			return "", err
		}
		err = isFolderInDirectory(serverPath, user.GetPath()) // Checks for folder's existence
		if err != nil {
			return "", err
		}
	}
	switch serverPath {
	case "\\", "/":
		path, err = user.setBackRoot()
	case "..":
		path, err = user.setBackDirectory()
	default:
		if filepath.IsAbs(serverPath) {
			path, err = user.setAbsDir(serverPath)
		} else {
			path, err = user.setForwardDir(serverPath)
		}
	}
	if err != nil {
		return path, err
	}
	return helper.GetVirtualStoragePath(path), nil
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
func (user *LoggedUser) RenameContent(contentPath string, newContentPath string) error {
	// Converting contentPath to absolute if it isn't
	if !filepath.IsAbs(contentPath) {
		err := IsFileInDirectory(contentPath, user.GetPath())
		if err != nil {
			return err
		}
		contentPath = helper.ConvertToAbsolute(user.GetPath(), contentPath) // if the file not abs -> file.* -> patn/file.*
	}
	return renameAbsContent(contentPath, newContentPath)
}

// Moves a content's path (files and folders)
func (user *LoggedUser) MoveContent(contentPath, newContentPath string) error {
	// Converting both contentPath and newContentPath to absolute paths if they aren't
	if !filepath.IsAbs(contentPath) {
		contentPath = helper.ConvertToAbsolute(user.GetPath(), contentPath)
	}
	if !filepath.IsAbs(newContentPath) {
		newContentPath = helper.ConvertToAbsolute(user.GetPath(), newContentPath)
	}
	return moveContent(contentPath, newContentPath)
}

// file operation that recieves all the file operations and send them to the function that is responsible for
func (user *LoggedUser) fileOperation(path string, operation func(string) error) error {
	if !filepath.IsAbs(path) {
		path = helper.ConvertToAbsolute(user.GetPath(), path) // convert to an absolute-server side path
	} else {
		if !strings.HasPrefix(path, helper.RootDir) {
			return &PremmisionError{path}
		}
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
	newPath := filepath.Dir(user.GetPath())
	if newPath == "." {
		return "", &PremmisionError{newPath}
	}
	err := user.SetPath(newPath)
	return user.GetPath(), err
}

// forward to next given dir, for example: cd homework11
func (user *LoggedUser) setForwardDir(forwardDir string) (string, error) {
	err := isFolderInDirectory(forwardDir, user.GetPath())
	if err != nil {
		return "", err
	}
	err = user.SetPath(filepath.Join(user.GetPath(), forwardDir))
	return user.GetPath(), err
}

func (user *LoggedUser) setAbsDir(absDir string) (string, error) {
	err := user.SetPath(absDir)
	return user.GetPath(), err
}

func (user *LoggedUser) ListContents() (string, error) {
	return getFolderContent(user.GetPath())
}

// create a file in the given directory
func createAbsFile(filePath string) error {
	err := validFileName(helper.Base(filePath), filepath.Dir(filePath)) // Validate file name
	if err != nil {
		return err
	}
	err = IsFileInDirectory(helper.Base(filePath), filepath.Dir(filePath))
	if err == nil { // If file exists
		return &FileExistError{helper.Base(filePath), filepath.Dir(filePath)}
	}
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
	err := validFileName(helper.Base(absDir), filepath.Dir(absDir)) // Validate folder name
	if err != nil {
		return err
	}
	err = isFolderInDirectory(helper.Base(absDir), filepath.Dir(absDir))
	if err == nil { // If folder exists
		return &FolderExistError{helper.Base(absDir), filepath.Dir(absDir)}
	}
	os.Mkdir(absDir, os.ModePerm) // Creates a folder
	return nil
}

// Remove an absolute file with the parameter file name
func removeAbsFile(filePath string) error {
	err := validFileName(helper.Base(filePath), filepath.Dir(filePath)) // Validate file name
	if err != nil {
		return err
	}
	err = IsFileInDirectory(helper.Base(filePath), filepath.Dir(filePath))
	if err != nil { // If file is not in directory
		return err
	}
	os.Remove(filePath)
	return nil
}

// Remove an absolute directory with the parameter folder name
func removeAbsFolder(absDir string) error {
	err := validFileName(helper.Base(absDir), filepath.Dir(absDir)) // Validate folder name
	if err != nil {
		return err
	}
	err = isFolderInDirectory(filepath.Base(absDir), filepath.Dir(absDir))
	if err != nil { // If folder is not in directory
		return err
	}
	os.RemoveAll(absDir)
	return nil
}

// Rename a file name in the current directory, gets full file path and new filename
func renameAbsContent(currentContentPath, newContentName string) error {
	newContentName = filepath.Join(filepath.Dir(currentContentPath), newContentName) // Get path
	// Use os.Rename to rename the file
	err := os.Rename(currentContentPath, newContentName)
	if err != nil {
		return err
	}
	return nil
}

func moveContent(currentAbsFilePath, newFileName string) error {
	return renameAbsContent(currentAbsFilePath, newFileName)
}

// Returns folder's content including its files and folders in a string variable. Return error if it fails
func getFolderContent(folderPath string) (string, error) {

	folder, err := os.Open(folderPath)
	if err != nil {
		return "", &OpenDirError{folderPath}
	}
	defer folder.Close()

	entries, err := folder.Readdir(-1) // Saves all entries in the directory

	if err != nil {
		return "", &ReadDirError{folderPath}
	}

	var builder strings.Builder
	var fileCounter, dirCounter uint

	for _, entry := range entries {
		// <DIR|FILE> <filename>
		if entry.IsDir() { // If the entry is a directory
			dirCounter++
			writeString(&builder, " <DIR> ", entry.Name(), "\n")
		} else { // Else if the entry is a file
			fileCounter++
			writeString(&builder, " <FILE> ", entry.Name(), "\n")
		}
	}
	writeString(&builder, fmt.Sprint((dirCounter)), "Dir(s), ", fmt.Sprint((fileCounter)), "File(s)")

	return builder.String(), nil
}

// Append entry to the builder.string object contains all the entries of the directory
func writeString(builder *strings.Builder, fileParameter ...string) {
	for _, stringObject := range fileParameter {
		builder.WriteString(stringObject)
	}
}
