package FileRequestsManager

import (
	"client/ClientErrors"
	"client/Helper"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type content struct {
	Name string `json:"name"` // File name (including its extension)
	Path string `json:"path"` // File's path in the Cloud
	Size uint32 `json:"size"` // File's size in bytes
}

// Creates a new file struct with the given parameters
func newContent(name string, path string, size uint32) content {
	return content{
		Name: name,
		Path: path,
		Size: size,
	}
}

// Checks local file and returns the file api if exists
func checkContent(filename string) (fs.FileInfo, error) {
	clearPath := strings.Replace(filename, "'", "", Helper.RemoveAll) // Clears ' encloused chars if they exist
	fileInfo, err := os.Stat(clearPath)                               // Check file (Access file path without enclosed quotation)
	if err != nil {
		if os.IsNotExist(err) { // If file not exists
			return nil, &ClientErrors.FileNotExistError{Filename: clearPath}
		} else {
			return nil, &ClientErrors.ReadFileInfoError{Filename: clearPath}
		}
	}
	return fileInfo, nil
}

// Returns directory size
func getDirSize(dirPath string) (uint32, error) {
	var totalSize int64
	// Walk through all the files and dirctories in the given dir to calculate its size
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil { // If couldn't read file info
			return &ClientErrors.ReadFileInfoError{Filename: info.Name()}
		}
		if !info.IsDir() {
			totalSize += info.Size() // Increase total size for files only
		}
		return nil
	})
	return uint32(totalSize), err
}
