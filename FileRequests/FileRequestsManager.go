package FileRequestsManager

import (
	"client/ClientErrors"
	"client/Helper"
	"client/Requests"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	commandArgumentIndex = 0
	pathArgumentIndex    = 0
	minimumArguments     = 1
	operationArguments   = 2
	oldFileName          = 0
	newFileName          = 1
	contentNameIndex     = 1

	rename_arguments = 2
	move_arguments   = 2

	showFolderArguments = 1
	path_index          = 1

	uploadFileArguments = 2
	uploadFileNameIndex = 0
	uploadFilePathIndex = 1

	CreateFileCommand   = "newfile"
	CreateFolderCommand = "newdir"

	RemoveFileCommand   = "rmfile"
	RemoveFolderCommand = "rmdir"
)

func convertResponeToPath(data string) string {
	parts := strings.Split(data, "CurrentDirectory:")
	return parts[path_index]
}

func HandleChangeDirectory(command_arguments []string, socket net.Conn) error {
	if len(command_arguments) < minimumArguments { // If argument was not given
		return &ClientErrors.InvalidArgumentCountError{Arguments: uint8(len(command_arguments)), Expected: uint8(operationArguments)}
	}

	data, err := Helper.ConvertStringToBytes(strings.Join(command_arguments, " "))
	if err != nil {
		return err
	}
	responeData, err := Requests.SendRequest(Requests.ChangeDirectoryRequest, data, socket)
	if err != nil {
		return err
	}

	path := convertResponeToPath(responeData)
	setCurrentPath(path)
	return nil
}

// Handle create content (file or directory) requests
func HandleCreate(command []string, socket net.Conn) error {
	if len(command) < operationArguments {
		return &ClientErrors.InvalidArgumentCountError{Arguments: uint8(len(command)), Expected: uint8(operationArguments)}
	}
	data, err := Helper.ConvertStringToBytes(strings.Join(command[contentNameIndex:], " "))
	if err != nil {
		return err
	}
	var createType Requests.RequestType

	switch command[commandArgumentIndex] {
	case CreateFileCommand:
		createType = Requests.CreateFileRequest
	case CreateFolderCommand:
		createType = Requests.CreateFolderRequest
	default:
		return fmt.Errorf("wrong create request")
	}
	_, err = Requests.SendRequest(createType, data, socket)
	return err
}

// Handle remove content (file or directory) requests
func HandleRemove(command []string, socket net.Conn) error {
	if len(command) != operationArguments {
		return &ClientErrors.InvalidArgumentCountError{Arguments: uint8(len(command)), Expected: uint8(operationArguments)}
	}
	data, err := Helper.ConvertStringToBytes(strings.Join(command[contentNameIndex:], " ")) // Convert content name to raw json bytes
	if err != nil {
		return err
	}

	var removeType Requests.RequestType

	switch command[commandArgumentIndex] {
	case RemoveFileCommand:
		removeType = Requests.DeleteFileRequest
	case RemoveFolderCommand:
		removeType = Requests.DeleteFolderRequest
	default:
		return fmt.Errorf("wrong remove request")
	}
	_, err = Requests.SendRequest(removeType, data, socket)
	return err
}

// Handle Rename request
func HandleRename(command_arguments []string, socket net.Conn) error {
	if len(command_arguments) < rename_arguments { // If argument was not given
		return &ClientErrors.InvalidArgumentCountError{Arguments: uint8(len(command_arguments)), Expected: uint8(rename_arguments)}
	}
	var oldcontentName string
	var newcontentName string
	if Helper.IsQuoted(command_arguments) { // Check if the command arguments are enclosed within a quotation (') marks
		oldcontentName = Helper.FindPath(command_arguments, Helper.FirstNameParameter)
		newcontentName = Helper.FindPath(command_arguments, Helper.SecondNameParameter)
	} else {
		oldcontentName = fmt.Sprintf("'" + command_arguments[oldFileName] + "'")
		newcontentName = fmt.Sprintf(" '" + command_arguments[newFileName] + "'")
	}
	paths := oldcontentName + newcontentName // Append to a string
	data, err := Helper.ConvertStringToBytes(paths)
	if err != nil {
		return err
	}
	_, err = Requests.SendRequest(Requests.RenameRequest, data, socket)
	return err
}

// Handle Move request
func HandleMove(command_arguments []string, socket net.Conn) error {
	if len(command_arguments) < move_arguments {
		return &ClientErrors.InvalidArgumentCountError{Arguments: uint8(len(command_arguments)), Expected: uint8(move_arguments)}
	}
	var currentFilePath string
	var newPath string
	if Helper.IsQuoted(command_arguments) { // Check if the command arguments are enclosed within a quotation (') marks
		currentFilePath = Helper.FindPath(command_arguments, Helper.FirstNameParameter) // Save the first path
		newPath = Helper.FindPath(command_arguments, Helper.SecondNameParameter)        // Save the second path
	} else {
		currentFilePath = fmt.Sprintf("'" + command_arguments[oldFileName] + "'")
		newPath = fmt.Sprintf(" '" + command_arguments[newFileName] + "'")
	}
	paths := currentFilePath + newPath // Appened to one string
	data, err := Helper.ConvertStringToBytes(paths)
	if err != nil {
		return err
	}
	_, err = Requests.SendRequest(Requests.MoveRequest, data, socket)
	return err
}

// Handle ls command (List contents command)
func HandleShow(command_arguments []string, socket net.Conn) (string, error) {
	if !(len(command_arguments) == showFolderArguments || len(command_arguments) == 0) { // check for amount of arguments
		return "", &ClientErrors.InvalidArgumentCountError{Arguments: uint8(len(command_arguments)), Expected: uint8(operationArguments)}
	}
	var data []byte
	var err error
	if len(command_arguments) == showFolderArguments { // If specific path has been specified
		data, err = Helper.ConvertStringToBytes(strings.Join(command_arguments[pathArgumentIndex:], " "))
		if err != nil {
			return "", err
		}
	} else {
		data = nil // If path hasn't been specified
	}
	respone, err := Requests.SendRequest(Requests.ShowRequest, data, socket)
	if err != nil {
		return "", err
	}
	return respone, nil
}

// Handle upload command
func HandleUpload(command_arguments []string, socket net.Conn) (string, error) {
	if len(command_arguments) < minimumArguments || len(command_arguments) > uploadFileArguments { // If file name was not provided or more parameters have been given
		return "", &ClientErrors.InvalidArgumentCountError{Arguments: uint8(len(command_arguments)), Expected: uint8(operationArguments)}
	}
	fileinfo, err := os.Stat(command_arguments[pathArgumentIndex])
	if err != nil {
		if os.IsNotExist(err) { // If file not exists
			return "", &ClientErrors.FileNotExistError{Filename: command_arguments[pathArgumentIndex]}
		} else {
			return "", &ClientErrors.ReadFileInfoError{Filename: command_arguments[pathArgumentIndex]}
		}
	}
	var path = ""                                      // Default path is empty (Relies on the server to pick the current directory as the path)
	if len(command_arguments) == uploadFileArguments { // If path has been specified
		path = command_arguments[uploadFilePathIndex]
	}
	fileSize := uint32(fileinfo.Size())
	file := newFile(command_arguments[uploadFileNameIndex], path, fileSize)

	//uploadFile almog.txt {optional: path}
}
