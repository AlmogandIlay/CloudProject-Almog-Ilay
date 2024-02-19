package FileSystem

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// File sizes
const (
	_  = iota
	KB = uint32(1 << (10 * iota)) // 1 KB = 1000 bytes
	MB = uint32(1 << (10 * iota)) // 1 MB = 1000000 bytes
	GB = uint32(1 << (10 * iota)) // 1 GB = 1000000000 bytes

	maxFileSize   = GB - 100*MB // 900MB Max file size
	maxNameLength = 256         // Max name length for Windows limitiations

	invalidFileCharacters = "\\/:*?\"<>|" // Invalid characters for Windows limitiations
)

type File struct {
	Name string // File name (including its extension)
	Path string // File's path in the Cloud
	Size uint32 // File's size in bytes
}

func newFile(name string, path string, size uint32) File {
	return File{
		Name: name,
		Path: path,
		Size: size,
	}
}

func (file *File) setPath(path string) {
	file.Path = path
}

// Validate file by its name and the file's size itself
func (user *LoggedUser) ValidateFile(file File) error {
	err := validFileName(file.Name, file.Path) // Check file name validation
	if err != nil {
		return err
	}
	err = validFileSize(file.Size) // Check file size validation
	if err != nil {
		return err
	}

	return nil
}

// Valids file size in terms of Cloud storage limitiations
func validFileSize(fileSize uint32) error {
	if fileSize > maxFileSize {
		return &FileSizeError{fileSize}
	}
	return nil
}

// Valids file size when creating a new file scenario
func (user *LoggedUser) validNewFileSize(fileSize uint32) error {
	rootSize, err := user.GetRootSize() // Get total amount of storage in usage
	if err != nil {
		return err
	}
	if fileSize+rootSize > maxFileSize { // If the new file size exceeds beyond the current storage space
		return &FileSizeError{fileSize}
	}
	return nil
}

// Checks if a file name is valid due to the Windows OS NTFS File System
func validFileName(name, path string) error {
	switch name {
	case "CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9": // Blacklist Windows filenames
		return &FileNameError{name}
	}
	if !validFileNameLength(name) {
		return &FileLengthError{name}
	}

	r := regexp.MustCompile(fmt.Sprintf("[%s]", regexp.QuoteMeta(invalidFileCharacters))) // regex check for specific file chars
	if r.MatchString(name) {                                                              // If filename has any chars from the blacklist file characters
		return &CharactersError{}
	}

	return nil
}

// Valids File Name Length due to Windows OS NTFS File System
func validFileNameLength(fileName string) bool {
	return uint32(len(fileName)) < maxNameLength && len(fileName) != 0
}

// Input: File Name (string), Path (string)
// Checks if the given file name exists in the given path. Returns error for invalid
// Output: error for file name that does not exist, any potential problems
func IsContentInDirectory(contentName, pathOfDir string) error {
	dir, err := os.Open(pathOfDir)

	if err != nil {
		return &OpenDirError{pathOfDir}
	}

	defer dir.Close()

	files, err := dir.ReadDir(-1) // Saves all the files in a slice
	if err != nil {
		return &ReadDirError{pathOfDir}
	}

	// check for files with same name while ignoring case sensitivity
	for _, file := range files {
		if strings.EqualFold(file.Name(), contentName) {
			return nil
		}
	}
	return &ContentNotExistError{contentName, pathOfDir}

}
