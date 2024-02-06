package FileSystem

import (
	"fmt"
	"strings"
)

type FileSizeError struct{ Size uint32 }
type FileNameError struct{ Name string }
type PathNotExistError struct{ Path string }

type OpenDirError struct{ Path string }
type ReadDirError struct{ Path string }
type FileExistError struct{ Name, Path string }
type FileLengthError struct{ Name string }
type CharactersError struct{}

type PremmisionError struct{ Path string }

// Creates custom errors for filesystem

func (fileError *FileSizeError) Error() string {
	return fmt.Sprintf("The file size is %d exceeded your total storage size which is %d", fileError.Size, -1)
}

func (fileError *FileNameError) Error() string {
	return fmt.Sprintf("The file: %s has invalid name", fileError.Name)
}

func (fileError *PathNotExistError) Error() string {
	return fmt.Sprintf("The path: %s not exist", fileError.Path)
}

func (fileError *OpenDirError) Error() string {
	return fmt.Sprintf("Cannot open %s dir", fileError.Path)
}

func (fileError *ReadDirError) Error() string {
	return fmt.Sprintf("Cannot read %s dir", fileError.Path)
}

func (fileError *FileExistError) Error() string {
	return fmt.Sprintf("the File %s already exist in %s path", fileError.Name, fileError.Path)
}

func (fileError *FileLengthError) Error() string {
	return fmt.Sprintf("The file %s is not the appropriate length, should be under %d", fileError.Name, maxFileSize)
}
func (fileError *CharactersError) Error() string {
	return fmt.Sprintf("Illegal letters such as: %s in the file", strings.Join(strings.Split(invalidFileCharacters, ""), " "))
}

func (fileError *PremmisionError) Error() string {
	return fmt.Sprintf("You have no permission to access %s", fileError.Path)
}
