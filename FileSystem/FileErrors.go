package FileSystem

import (
	helper "CloudDrive/Helper"
	"fmt"
	"strings"
)

type FileSizeError struct{ Size uint32 }
type FileNameError struct{ Name string }
type PathNotExistError struct{ Path string }

type OpenDirError struct{ Path string }
type ReadDirError struct{ Path string }
type FileExistError struct{ Name, Path string }
type FolderExistError struct{ Name, Path string }
type FileNotExistError struct{ Name, Path string }
type FolderNotExistError struct{ Name, Path string }
type ContentNotExistError struct{ Name, Path string }
type ContentExistError struct{ Name, Path string }
type FileLengthError struct{ Name string }
type CharactersError struct{}
type SizeCalculationError struct{}

type PremmisionError struct{ Path string }

type InitializeError struct{}

type UnmarshalError struct{}

// Creates custom errors for filesystem

func (fileError *FileSizeError) Error() string {
	return fmt.Sprintf("The file size %d has exceeded your total storage size which is %d", fileError.Size, -1)
}

func (fileError *FileNameError) Error() string {
	return fmt.Sprintf("The file: '%s' has invalid name", fileError.Name)
}

func (fileError *PathNotExistError) Error() string {
	return fmt.Sprintf("The path: '%s' not exist", helper.GetVirtualStoragePath(fileError.Path))
}

func (fileError *OpenDirError) Error() string {
	return fmt.Sprintf("Cannot open '%s' dir", helper.GetVirtualStoragePath(fileError.Path))
}

func (fileError *ReadDirError) Error() string {
	return fmt.Sprintf("Cannot read '%s' dir", helper.GetVirtualStoragePath(fileError.Path))
}

func (fileError *FileExistError) Error() string {
	return fmt.Sprintf("The file '%s' is already exist in '%s'", fileError.Name, helper.GetVirtualStoragePath(fileError.Path))
}

func (fileError *FolderExistError) Error() string {
	return fmt.Sprintf("The folder '%s' already exist in path '%s'", fileError.Name, helper.GetVirtualStoragePath(fileError.Path))
}

func (fileError *FileNotExistError) Error() string {
	return fmt.Sprintf("The file '%s' not exist in '%s' path", fileError.Name, helper.GetVirtualStoragePath(fileError.Path))
}

func (fileError *FolderNotExistError) Error() string {
	return fmt.Sprintf("The folder '%s' not exist in '%s' path", fileError.Name, helper.GetVirtualStoragePath(fileError.Path))
}

func (fileError *ContentNotExistError) Error() string {
	return fmt.Sprintf("The content '%s' does not exist in '%s' path", fileError.Name, helper.GetVirtualStoragePath(fileError.Path))
}

func (fileError *ContentExistError) Error() string {
	return fmt.Sprintf("The content '%s' is already exists in %s path", fileError.Name, helper.GetVirtualStoragePath(fileError.Path))
}

func (fileError *FileLengthError) Error() string {
	return fmt.Sprintf("The file '%s' is not the appropriate length, should be under '%d'", fileError.Name, maxFileSize)
}
func (fileError *CharactersError) Error() string {
	return fmt.Sprintf("Cannot use Illegal letters such as: '%s' in the name", strings.Join(strings.Split(invalidFileCharacters, ""), " "))
}
func (fileError *SizeCalculationError) Error() string {
	return "There has been an error when calculating the root path"
}

func (fileError *PremmisionError) Error() string {
	return fmt.Sprintf("You have no permission to access out of your root%s", helper.GetVirtualStoragePath(fileError.Path))
}

func (fileError *InitializeError) Error() string {
	return "There has been an error when attempting to access the allocated memory"
}

func (fileError *UnmarshalError) Error() string {
	return "There has been an error when attempting to encode the requested file."
}
