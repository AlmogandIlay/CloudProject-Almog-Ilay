package FileSystem

import (
	"os"
	"regexp"
	"strings"
)

const (
	_  = iota
	KB = uint32(1 << (10 * iota)) // 1 << (10 * 1)
	MB = uint32(1 << (10 * iota))
	GB = uint32(1 << (10 * iota))

	maxFileSize   = GB - 100*MB // 900KB
	maxNameLength = 256

	invalidFileCharacters = "\\/:*?\"<>|"
)

type File struct {
	Name string // name  + "." + extension
	Path string // path file in the cloud/server
	Size uint32 // in bytes
}

func (user *LoggedUser) NewFile(name string, size uint32) (*File, error) {

	err := ValidFileName(name, user.CurrentPath)
	if err != nil {
		return nil, err
	}
	err = ValidFileSize(size)
	if err != nil {
		return nil, err
	}
	return &File{name, user.CurrentPath, size}, nil
}

// *instead check for size < freeUserMemory
func ValidFileSize(fileSize uint32) error {
	if fileSize > maxFileSize {
		return &FileSizeError{fileSize}
	}
	return nil
}

// check for file with same name in folder
func ValidFileName(name, path string) error {
	switch name {
	case "CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9": // Blacklist filenames
		return &FileNameError{name}
	}
	if !validFileNameLength(name) {
		return &FileLengthError{name}
	}
	if regexp.MustCompile(invalidFileCharacters).MatchString(string(name)) {
		return &CharactersError{}
	}
	return isFileInDirectory(name, path)
}

func validFileNameLength(fileName string) bool {
	return uint32(len(fileName)) < maxNameLength && len(fileName) != 0
}

// transfer to helper

func isFileInDirectory(fileName, pathOfDir string) error {
	dir, err := os.Open(pathOfDir)

	if err != nil {
		return &OpenDirError{pathOfDir}
	}

	defer dir.Close()

	files, err := dir.ReadDir(-1) // Saves all the path files in a slice
	if err != nil {
		return &ReadDirError{pathOfDir}
	}

	// check for file with same name. note: files are case-senstive
	for _, file := range files {
		if strings.EqualFold(file.Name(), fileName) {
			return &FileExistError{fileName, pathOfDir}
		}
	}

	return nil
}
