package Helper

import (
	"os"
	"path/filepath"
)

// Returns whether the given path exists. Returns error if the check gone wrong
func IsPathExists(path string) (bool, error) {
	if path != "" { // If path has any value, otherwise returns true automatically
		_, err := os.Stat(path)
		if err == nil { // If path exists
			return true, nil
		}
		if os.IsNotExist(err) { // If path not exists
			return false, nil
		}
		// If the check gone wrong
		return false, err
	}
	return true, nil
}

// Counting contents (files and folders) in a given folderpath
func CountContents(folderPath string) (uint, uint, error) {
	var filesCount, folderCount uint

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		// for each content in the directory
		if err != nil {
			return err
		}
		if info.IsDir() { // If it's a directory
			folderCount++
		} else { // If it's a file
			filesCount++
		}

		return nil
	})
	// Substract 1 from foldercount to remove the counting of the folderPath itself
	return filesCount, folderCount - 1, err
}
