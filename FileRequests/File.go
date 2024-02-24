package FileRequestsManager

import (
	"client/ClientErrors"
	"client/Helper"
	"io/fs"
	"os"
	"strings"
)

type file struct {
	Name string `json:"name"` // File name (including its extension)
	Path string `json:"path"` // File's path in the Cloud
	Size uint32 `json:"size"` // File's size in bytes
}

// Creates a new file struct with the given parameters
func newFile(name string, path string, size uint32) file {
	return file{
		Name: name,
		Path: path,
		Size: size,
	}
}

// Checks local file and returns the file api if exists
func checkFile(filename string) (fs.FileInfo, error) {
	clearPath := strings.Replace(filename, "'", "", Helper.RemoveAll)
	fileInfo, err := os.Stat(clearPath) // Check file (Access file path without enclosed quotation)
	if err != nil {
		if os.IsNotExist(err) { // If file not exists
			return nil, &ClientErrors.FileNotExistError{Filename: clearPath}
		} else {
			return nil, &ClientErrors.ReadFileInfoError{Filename: clearPath}
		}
	}
	return fileInfo, nil
}
