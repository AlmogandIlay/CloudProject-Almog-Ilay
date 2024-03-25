package FileSystem

import (
	helper "CloudDrive/Helper"
	"fmt"
	"strings"
)

type FileSizeError struct{ Size uint32 }
type ContentNameError struct{ Name string }
type PathNotExistError struct{ Path string }

type OpenDirError struct{ Path string }
type ReadDirError struct{ Path string }
type AbsFileError struct{ Path string }
type RelFileError struct{ Path string }
type FileExistError struct{ Name, Path string }
type FileInfoError struct{ Name string }
type FolderExistError struct{ Name, Path string }
type FileNotExistError struct{ Name, Path string }
type RareIssueWithFile struct{ Name string }
type RenameError struct{ Name, NewName string }
type FolderNotExistError struct{ Name, Path string }
type ContentNotExistError struct{ Name, Path string }
type ContentExistError struct{ Name, Path string }
type ContentLengthError struct{ Name string }
type CharactersError struct{}
type SizeCalculationError struct{}
type FileExceededCurrentAvailableStorage struct{ Name string }

type PremmisionOutOfRootError struct{}
type PremmisionError struct{ Path string }

type InitializeError struct{}

type UnmarshalError struct{}
type MarshalError struct{}

type CreatePrivateSocketError struct{}

type UploadTimeOut struct{}

// Creates custom errors for filesystem

func (fileError *FileSizeError) Error() string {
	return fmt.Sprintf("The file size %d has exceeded your total storage size which is %d", fileError.Size, -1)
}

func (fileError *ContentNameError) Error() string {
	return fmt.Sprintf("The content name '%s' is invalid", fileError.Name)
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

func (fileError *FileInfoError) Error() string {
	return fmt.Sprintf("Couldn't read filename '%s''s info", fileError.Name)
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
	return fmt.Sprintf("The content '%s' is already exists in '%s' path", fileError.Name, helper.GetVirtualStoragePath(fileError.Path))
}

func (fileError *RareIssueWithFile) Error() string {
	return fmt.Sprintf("There has been a rare issue with the provided file: '%s' when checking its size.\nPlease report that to the developers.", fileError.Name)
}

func (fileError *RenameError) Error() string {
	return fmt.Sprintf("Connot rename the file: %s to %s", fileError.Name, fileError.NewName)
}

func (fileError *ContentLengthError) Error() string {
	return fmt.Sprintf("The content '%s''s length is not an appropriate length", fileError.Name, maxcontentSize)
}
func (fileError *CharactersError) Error() string {
	return fmt.Sprintf("Cannot use Illegal letters such as: '%s' in the name", strings.Join(strings.Split(invalidContentCharacters, ""), " "))
}
func (fileError *SizeCalculationError) Error() string {
	return "There has been an error when calculating the root path"
}

func (fileError *FileExceededCurrentAvailableStorage) Error() string {
	return fmt.Sprintf("The file '%s' has exceeded your current available storage.\nPlease clean your storage", fileError.Name)
}

func (fileError *PremmisionError) Error() string {
	return fmt.Sprintf("You have no permission to access the path: %s", helper.GetVirtualStoragePath(fileError.Path))
}

func (fileError *PremmisionOutOfRootError) Error() string {
	return "You have no permission to access out of your root"
}

func (fileError *InitializeError) Error() string {
	return "There has been an error when attempting to access the allocated memory"
}

func (fileError *UnmarshalError) Error() string {
	return "There has been an error when attempting to decode the requested file."
}

func (fileError *MarshalError) Error() string {
	return "There has been an error when attempting to encode the requested file."
}

func (fileError *CreatePrivateSocketError) Error() string {
	return "Couldn't create a private socket for the file upload."
}

func (fileError *UploadTimeOut) Error() string {
	return "Upload process has been stopped as it's passed the timeout duration of upload.\nThis doesn't mean the upload proccess didn't finished.\nPlease take a look of the uploaded content."
}

func (fileError *AbsFileError) Error() string {
	return fmt.Sprintf("The argument %s should be pure name of file, not abs path", fileError.Path)
}

func (fileError *RelFileError) Error() string {
	return fmt.Sprintf("Couldn't convert '%s' path to Relative path", fileError.Path)
}
