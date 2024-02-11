package FileSystem

import (
	"os"
)

type Folder struct {
	Name    string    // File name without extension
	Path    string    // File path
	Files   []File    // Files contain in the folder
	Folders []*Folder // Slice of pointers to Folder structs
}

func isFolderInDirectory(path, pathOfDir string) error {
	user_path := pathOfDir + "\\" + path
	_, err := os.Open(user_path)

	if err != nil {
		return &PathNotExistError{Path: path}
	}
	return nil
}
