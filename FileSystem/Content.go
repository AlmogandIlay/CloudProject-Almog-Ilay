package FileSystem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// content sizes
const (
	_  = iota
	KB = uint32(1 << (10 * iota)) // 1 KB = 1000 bytes
	MB = uint32(1 << (10 * iota)) // 1 MB = 1000000 bytes
	GB = uint32(1 << (10 * iota)) // 1 GB = 1000000000 bytes

	maxcontentSize = GB - 100*MB // 900MB Max content size
	maxNameLength  = 256         // Max name length for Windows limitiations

	invalidContentCharacters = "\\/:*?\"<>|" // Invalid characters for Windows limitiations

	nonSize = 0
)

type Content struct {
	Name string // content name (including its extension)
	Path string // content's path in the Cloud
	Size uint32 // content's size in bytes
}

func newContent(name string, path string, size uint32) Content {
	return Content{
		Name: name,
		Path: path,
		Size: size,
	}
}

// Validate content by its name and the content's size itself
func (user *LoggedUser) ValidateContent(content Content) error {
	err := validContentName(content.Name) // Check Content name validation
	if err != nil {
		return err
	}
	err = validContentSize(content.Size) // Check Content size validation
	if err != nil {
		return err
	}

	return nil
}

// Valids content size in terms of Cloud storage limitiations
func validContentSize(contentSize uint32) error {
	if contentSize > maxcontentSize {
		return &FileSizeError{contentSize}
	}
	return nil
}

// Valids content size when creating a new content scenario
func (user *LoggedUser) validNewContentSize(contentSize uint32, contentName string) error {
	rootSize, err := user.GetRootSize() // Get total amount of storage in usage
	if err != nil {
		return err
	}
	if contentSize+rootSize > maxcontentSize { // If the new content size exceeds beyond the current storage space
		return &FileExceededCurrentAvailableStorage{Name: contentName}
	}
	return nil
}

// Checks if a content name is valid due to the Windows OS NTFS content System
func validContentName(name string) error {
	switch name {
	case "CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9": // Blacklist Windows filenames
		return &FileNameError{name}
	}
	if !validContentNameLength(name) {
		return &FileLengthError{name}
	}

	r := regexp.MustCompile(fmt.Sprintf("[%s]", regexp.QuoteMeta(invalidContentCharacters))) // regex check for specific content chars
	if r.MatchString(name) {                                                                 // If contentname has any chars from the blacklist content characters
		return &CharactersError{}
	}

	return nil
}

// Valids content Name Length due to Windows OS NTFS content System
func validContentNameLength(contentName string) bool {
	return uint32(len(contentName)) < maxNameLength && len(contentName) != 0
}

// Input: content Name (string), Path (string)
// Checks if the given content name exists in the given path. Returns error for invalid
// Output: error for content name that does not exist, any potential problems
func IsContentInDirectory(contentName, pathOfDir string) error {
	dir, err := os.Open(pathOfDir)

	if err != nil {
		return &OpenDirError{pathOfDir}
	}

	defer dir.Close()

	files, err := dir.ReadDir(-1) // Saves all the contents in a slice
	if err != nil {
		return &ReadDirError{pathOfDir}
	}

	// check for contents with same name while ignoring case sensitivity
	for _, file := range files {
		if strings.EqualFold(file.Name(), contentName) {
			return nil
		}
	}
	return &ContentNotExistError{contentName, pathOfDir}

}

// Returns file's size by its absolute filename
func getFileSize(absFilename string) (uint32, error) {
	fileInfo, err := os.Stat(absFilename)
	if err != nil {
		if os.IsNotExist(err) {
			// If file not exists
			return nonSize, &FileNotExistError{Name: filepath.Base(absFilename), Path: filepath.Dir(absFilename)}
		}
		return nonSize, &RareIssueWithFile{Name: filepath.Base(absFilename)} // If unlikely and very rare error has happened report to client without leaking the details

	}
	fileSize := uint32(fileInfo.Size()) // Get file's size

	return fileSize, nil
}

// Parse json data request to content struct
func ParseDataToContent(data json.RawMessage) (*Content, error) {
	var content Content

	err := json.Unmarshal(data, &content) // Convert json request to content struct
	if err != nil {                       // If conversion failed
		return nil, &UnmarshalError{} // Convert the error to our custom made error.
	}
	return &content, nil

}
