package FileSystem

import (
	helper "CloudDrive/Helper"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	emptyChunks  = 0
	emptySize    = 0
	errorRespone = 999
	validRespone = 200
)

// Change the user current directory into garbage directory
func (user *LoggedUser) GarbageChangeDirectory() (string, error) {
	return user.ChangeDirectory(helper.GetGarbagePath(user.UserID))
}

// FileManager: API interface for LoggedUser to interact with the file commands in the cloud drive.

// Removes back subdriectories (..) from the absolute client-side path
func removeBackDirectories(absPath string) (string, error) {
	// Convert the absolute-server side path to absolute-client side path for easy access in the split proccess
	absPath = helper.GetVirtualStoragePath(absPath)
	subDirectories := strings.Split(absPath, "\\")
	subDirectories[0] = helper.RootDir // Manually modify the first subDirectory which is always 'Root:\' due to the split with '\'
	var newDirectories []string        // The actual subDirectories the client meant to access
	for i := 0; i < len(subDirectories); i++ {
		if subDirectories[i] != ".." { // If subDirectory is not back directory
			newDirectories = append(newDirectories, subDirectories[i])
		} else {
			// If subDirectory is a back directory (..)

			// Checks if not trying to go back from the Root:\ (avoid attempts to escape the dedicated memory space)
			if newDirectories[len(newDirectories)-1] == helper.RootDir {
				return "", &PremmisionOutOfRootError{}
			}
			newDirectories = newDirectories[:len(newDirectories)-1] // Removes last subdirectory of the actual subdriectory
		}

	}
	return strings.Join(newDirectories, "\\"), nil
}

// Changes the current directory for the user according to the parameter
func (user *LoggedUser) ChangeDirectory(parameter string) (string, error) {
	var err error
	var path string
	serverPath := helper.GetServerStoragePath(user.UserID, parameter)  // Convert client paramer path to server side if it's absolute or has sub directories
	if serverPath != ".." && serverPath != "\\" && serverPath != "/" { // If the path is not going back or forward
		err = validContentName(helper.Base(serverPath)) // Valid for files and folders are equals. calling the validFileName
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

		// If path is absolute (starts with DrivePath)
		if filepath.IsAbs(serverPath) {

			if strings.Contains(serverPath, "\\..") { // If absolute path contains back directory (..)
				// Removes back directories from the path and return the actual path the user meant to go into
				serverPath, err = removeBackDirectories(serverPath)               // Returns path without back directories as absolute client path
				serverPath = helper.GetServerStoragePath(user.UserID, serverPath) // Convert back the absolute client path to absolute server path
			}
			// If client tried to access content in Garbage path
			if helper.IsContainGarbage(serverPath, user.UserID) && len(serverPath) != len(helper.GetGarbagePath(user.UserID)) {
				return "", &PremmisionError{serverPath}
			}
			if err != nil { // check error for removeBackDirectories but first check if it contains in Garbage as a priotity error check
				return "", err
			}
			path, err = user.setAbsDir(serverPath)

		} else {
			if helper.IsContainGarbage(helper.ConvertToAbsolute(user.CurrentPath, serverPath), user.UserID) && len(helper.ConvertToAbsolute(user.CurrentPath, serverPath)) != len(helper.GetGarbagePath(user.UserID)) {
				return "", &PremmisionError{helper.ConvertToAbsolute(user.CurrentPath, serverPath)}
			}
			path, err = user.setForwardDir(serverPath)
		}
	}
	if err != nil {
		return path, err
	}
	return helper.GetVirtualStoragePath(path), nil
}

// Creates a new file in the given path
func (user *LoggedUser) CreateFile(fileName string) error {

	// any path into server absulote path
	serverpath := helper.GetServerStoragePath(user.UserID, fileName)
	if !filepath.IsAbs(serverpath) {
		serverpath = helper.ConvertToAbsolute(user.GetPath(), serverpath) // convert to an absolute-server side path
	}

	if helper.IsContainGarbage(serverpath, user.UserID) {
		return &PremmisionError{helper.GetGarbagePath(user.UserID)}
	}
	return user.fileOperation(serverpath, createAbsFile)
}

// Creates a new folder in the current directory
func (user *LoggedUser) CreateFolder(folderName string) error {
	// any path into server absulote path
	serverpath := helper.GetServerStoragePath(user.UserID, folderName) // Convert folderName client path side to absolute if given absolute
	if !filepath.IsAbs(serverpath) {                                   // Convert path folder name to absolute if it doesn't
		serverpath = helper.ConvertToAbsolute(user.GetPath(), serverpath) // convert to an absolute-server side path
	}
	if helper.IsContainGarbage(serverpath, user.UserID) { // If the folder to be created should be in Garbage folder
		return &PremmisionError{helper.GetGarbagePath(user.UserID)}
	}
	return user.fileOperation(serverpath, createAbsDir)
}

// Removes a content
func (user *LoggedUser) RemoveContent(contentName string) error {
	// any path into server absulote path
	serverpath := helper.GetServerStoragePath(user.UserID, contentName) // Convert client side path to server side (convert to absolute only if contentName is absolute)
	if !filepath.IsAbs(serverpath) {
		serverpath = helper.ConvertToAbsolute(user.GetPath(), serverpath) // convert to an absolute-server side path if it doesn't
	}

	err := IsContentInDirectory(filepath.Base(serverpath), filepath.Dir(serverpath)) // Checks if the content to remove does exists in its path
	if err != nil {                                                                  // If content doesn't exist in its path
		return &ContentNotExistError{Name: filepath.Base(serverpath), Path: filepath.Dir(serverpath)}
	}

	if helper.IsContainGarbage(serverpath, user.UserID) { // If the content path is in the Garbage path
		return user.fileOperation(serverpath, removeAbsContent) // Delete the content
	} else { // If the content doesn't in the Garbage path
		newContentPath := filepath.Join(helper.GetGarbagePath(user.UserID), filepath.Base(serverpath)) // Rename destination path for moving content
		err = IsContentInDirectory(filepath.Base(serverpath), filepath.Dir(newContentPath))
		if err == nil { // Plaster to solve sending to garbage a filename that already exists
			return &ContentExistError{Name: filepath.Base(serverpath), Path: filepath.Dir(newContentPath)}
		}

		moveContent(serverpath, newContentPath) // Move content to Garbage path
		return nil
	}
}

// Renames a content
func (user *LoggedUser) RenameContent(contentPath string, newContentPath string) error {
	// Converting contentPath to absolute if it doesn't
	if !helper.IsAbs(contentPath) { // Converting the content to an absolute path if it doesn't
		contentPath = helper.ConvertToAbsolute(user.GetPath(), contentPath) // if the file not abs -> file.* -> patn/file.*
	}
	contentPath = helper.GetServerStoragePath(user.UserID, contentPath)                // Convert absolute-client side path to absolute-server side path
	err := IsContentInDirectory(filepath.Base(contentPath), filepath.Dir(contentPath)) // Check if the content that will be renamed does exist
	if err != nil {
		return err
	}
	err = IsContentInDirectory(filepath.Base(newContentPath), filepath.Dir(contentPath)) // Check if the new content name does not exist already
	if err == nil {
		return &ContentExistError{Name: filepath.Base(newContentPath), Path: filepath.Dir(contentPath)}
	}
	err = user.ValidateContent(newContent(filepath.Base(newContentPath), filepath.Dir(contentPath), 0)) // Validate the newcontentPath name
	if err != nil {
		return err
	}

	garbagePath := helper.GetGarbagePath(user.UserID)
	// check if the path starts with the garbage path
	if strings.HasPrefix(helper.GetServerStoragePath(user.UserID, contentPath), garbagePath) {
		return &PremmisionError{helper.GetServerStoragePath(user.UserID, contentPath)}
	}
	renameAbsContent(contentPath, newContentPath)
	return nil
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

	err := IsContentInDirectory(filepath.Base(contentPath), filepath.Dir(contentPath)) // Checks if the content to move does exists in its path
	if err != nil {
		return err
	}
	err = IsContentInDirectory(filepath.Base(newContentPath), filepath.Dir(newContentPath)) // Checks if the to move request provided directory exists
	if err != nil {
		return err
	}

	newContentPath += "\\" + filepath.Base(contentPath) // Add the content (file/path)'s name extension to the new directory
	moveContent(contentPath, newContentPath)
	return nil
}

// list all the files in the given (or not given) path
func (user *LoggedUser) ListContents(path string) (string, error) {
	if path == "" { // If path hasn't been specified
		path = user.GetPath() // put current directory as default
	} else { // If path has been specified
		if !helper.IsAbs(path) { // Convert the path to absolute if it doesn't
			path = helper.ConvertToAbsolute(user.GetPath(), path)
		}
		path = helper.GetServerStoragePath(user.UserID, path) // Convert absolute-client path to absolute-server path
	}

	if helper.IsContainGarbage(path, user.UserID) && len(path) != len(helper.GetGarbagePath(user.UserID)) {
		return "", &PremmisionError{path}
	}
	return getFolderContent(path)
}

// file operation that recieves all the file operations and send them to the function that is responsible for
func (user *LoggedUser) fileOperation(path string, operation func(string) error) error {

	if !strings.HasPrefix(path, helper.DrivePath) {
		return &PremmisionOutOfRootError{}
	}

	return operation(helper.GetServerStoragePath(user.UserID, path)) // Sends the real path file argument
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
		return "", &PremmisionError{}
	}
	err := user.SetPath(newPath)
	return user.GetPath(), err
}

// forward to next given dir, for example: cd homework11
func (user *LoggedUser) setForwardDir(forwardDir string) (string, error) {
	err := user.SetPath(filepath.Join(user.GetPath(), forwardDir))
	return user.GetPath(), err
}

func (user *LoggedUser) setAbsDir(absDir string) (string, error) {
	err := user.SetPath(absDir)
	return user.GetPath(), err
}

// create a file in the given directory
func createAbsFile(filePath string) error {
	err := validContentName(helper.Base(filePath)) // Validate file name
	if err != nil {
		return err
	}
	err = IsContentInDirectory(helper.Base(filePath), filepath.Dir(filePath))
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
	err := validContentName(helper.Base(absDir)) // Validate folder name
	if err != nil {
		return err
	}
	err = IsContentInDirectory(helper.Base(absDir), filepath.Dir(absDir))
	if err == nil { // If folder exists
		return &FolderExistError{helper.Base(absDir), filepath.Dir(absDir)}
	}
	os.Mkdir(absDir, os.ModePerm) // Creates a folder
	return nil
}

// Remove a content with the given absolute contentPath parameter
func removeAbsContent(contentPath string) error {
	contentInfo, err := os.Stat(contentPath) // Get the content info
	if err != nil {                          // If gathering info has failed which should be a rare problem, report that
		return &RareIssueWithFile{Name: helper.Base(contentPath)}
	}
	switch mode := contentInfo.Mode(); { // Checks content's type
	case mode.IsDir(): // If content is Directory
		os.RemoveAll(contentPath) // Remove directory

	case mode.IsRegular(): // If content is File
		os.Remove(contentPath) // Remove file
	}

	return nil
}

// Renames a content
func renameAbsContent(currentContentPath, newContentName string) {
	if !filepath.IsAbs(newContentName) { // If the file to rename isn't absolute
		newContentName = filepath.Join(filepath.Dir(currentContentPath), newContentName) // Saves the new filename as absolute
	}
	// Use os.Rename to rename the file
	err := os.Rename(currentContentPath, newContentName)
	if err != nil { // If something bad happened, crash and report it in the server instead of reporting it to the client
		log.Fatal(err.Error())
	}
}

func moveContent(currentAbsFilePath, newFileName string) {
	renameAbsContent(currentAbsFilePath, newFileName)
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
	writeString(&builder, fmt.Sprint((dirCounter)), " Dir(s), ", fmt.Sprint((fileCounter)), " File(s)")

	return builder.String(), nil
}

// Append entry to the builder.string object contains all the entries of the directory
func writeString(builder *strings.Builder, fileParameter ...string) {
	for _, stringObject := range fileParameter {
		builder.WriteString(stringObject)
	}
}

// Duplicate from ResponeInfo to avoid import cycle. build an error responeInfo struct.
func buildError(response string) interface{} {
	return struct {
		Type    int    `json:"Type"`
		Respone string `json:"Data"`
	}{
		Type:    errorRespone,
		Respone: response,
	}
}

func buildRespone(responeType int, response []byte) interface{} {
	return struct {
		Type    int    `json:"Type"`
		Respone string `json:"Data"`
	}{
		Type:    responeType,
		Respone: string(response),
	}
}

// Duplicate from ResponeInfo to avoid import cycle. send Respone info.
func sendResponseInfo(conn *net.Conn, responseInfo interface{}) error {
	message, _ := json.Marshal(responseInfo)
	return helper.SendData(conn, message)
}
