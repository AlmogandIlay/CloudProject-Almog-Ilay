package FileRequestsManager

import (
	"client/ClientErrors"
	"client/Helper"
	"client/Requests"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"strings"
)

const (
	// Command Indexes:
	/////////////////////////
	commandArgumentIndex     = 0
	pathArgumentIndex        = 0
	minimumArguments         = 1
	operationArguments       = 2
	oldFileName              = 0
	newFileName              = 1
	contentNameIndex         = 1
	rename_arguments         = 2
	move_arguments           = 2
	showFolderArguments      = 1
	minimumdownloadArguments = 1
	localPathIndex           = 2
	cloudPathIndex           = 2
	/////////////////////////

	path_index = 1

	// Commands:
	CreateFileCommand   = "newfile"
	CreateFolderCommand = "newdir"
	RemoveFileCommand   = "rmfile"
	RemoveFolderCommand = "rmdir"
)

func convertResponeToPath(data string) string {
	parts := strings.Split(data, "CurrentDirectory:")
	return parts[path_index]
}

func HandleChangeDirectory(command_arguments []string, socket *net.Conn) error {
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

// Handle Garbage request
func HandleGarbage(socket *net.Conn) error {
	responeData, err := Requests.SendRequest(Requests.GarbageRequest, nil, socket) // Send request type without any data

	if err != nil {
		return err
	}

	path := convertResponeToPath(responeData)
	setCurrentPath(path)
	return nil
}

// Handle create content (file or directory) requests
func HandleCreate(command []string, socket *net.Conn) error {
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
func HandleRemove(command []string, socket *net.Conn) error {
	if len(command) < operationArguments {
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
func HandleRename(command_arguments []string, socket *net.Conn) error {
	if len(command_arguments) < rename_arguments { // If argument was not given
		return &ClientErrors.InvalidArgumentCountError{Arguments: uint8(len(command_arguments)), Expected: uint8(rename_arguments)}
	}
	var oldcontentName string
	var newcontentName string
	if Helper.IsQuoted(command_arguments, Helper.TwoCloudPaths) { // Check if the command arguments are enclosed within a quotation (') marks
		oldcontentName = Helper.FindPath(command_arguments, Helper.FirstNameParameter, Helper.TwoCloudPaths)
		newcontentName = fmt.Sprintf(" '" + Helper.FindPath(command_arguments, Helper.SecondNameParameter, Helper.TwoCloudPaths) + "'")
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
func HandleMove(command_arguments []string, socket *net.Conn) error {
	if len(command_arguments) < move_arguments {
		return &ClientErrors.InvalidArgumentCountError{Arguments: uint8(len(command_arguments)), Expected: uint8(move_arguments)}
	}
	var currentFilePath string
	var newPath string
	if Helper.IsQuoted(command_arguments, Helper.TwoCloudPaths) { // Check if the command arguments are enclosed within a quotation (') marks
		currentFilePath = Helper.FindPath(command_arguments, Helper.FirstNameParameter, Helper.TwoCloudPaths)                    // Save the first path
		newPath = fmt.Sprintf(" '" + Helper.FindPath(command_arguments, Helper.SecondNameParameter, Helper.TwoCloudPaths) + "'") // Save the second path
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
func HandleShow(command_arguments []string, socket *net.Conn) (string, error) {
	if !(len(command_arguments) >= showFolderArguments || len(command_arguments) == 0) { // check for amount of arguments
		return "", &ClientErrors.InvalidArgumentCountError{Arguments: uint8(len(command_arguments)), Expected: uint8(operationArguments)}
	}
	var data []byte
	var err error
	if len(command_arguments) >= showFolderArguments { // If specific path has been specified
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

// Handles upload command
func HandleUploadFile(command_arguments []string, socket *net.Conn) error {
	if len(command_arguments) < minimumArguments { // If file name was not provided
		return &ClientErrors.InvalidArgumentCountError{Arguments: uint8(len(command_arguments)), Expected: uint8(operationArguments)}
	}
	var filename string
	var cloudpath string
	if Helper.IsQuoted(command_arguments, Helper.OneClosedPath) { // Check if the first command argument is enclosed within a quotation (') marks
		filename = Helper.FindPath(command_arguments, Helper.FirstNameParameter, Helper.OneClosedPath) // Save the first path (filename to upload)
		if Helper.IsQuoted(command_arguments, Helper.TwoCloudPaths) {                                  // Check if the second command argument is enclosed within a quotation (') marks
			cloudpath = Helper.FindPath(command_arguments, Helper.SecondNameParameter, Helper.TwoCloudPaths) // Save the second path (path in cloud storage to save)
		} else { // If first path is quoted but the second doesn't
			cloudpath = Helper.ReturnNonQuotedSecondPath(command_arguments)
		}
	} else { // If command arguments are not enclosed within a quotation (') marks
		// relay on argument indexes
		filename = command_arguments[oldFileName]
		if len(command_arguments) == cloudPathIndex { // If client specificed a path to save in cloud
			cloudpath = command_arguments[newFileName]
		}
	}

	fileInfo, err := checkFile(filename) // Check if file exists, if it does returns file info api
	if err != nil {
		return err
	}
	fileSize := uint32(fileInfo.Size())
	file := newFile(filepath.Base(strings.Replace(filename, "'", "", Helper.RemoveAll)), cloudpath, fileSize) // Create a new file struct for server communication
	file_data, err := json.Marshal(file)
	if err != nil {
		return &ClientErrors.JsonEncodeError{}
	}

	respone, err := Requests.SendRequest(Requests.UploadFileRequest, file_data, socket) // Sends upload file request
	if err != nil {                                                                     // If upload file request was rejected
		return err
	}
	chunksSize, err := Helper.ConvertResponeToChunks(respone) // Convert respone to chunks size
	if err != nil {                                           // If chunks size was returned from the server in a wrong type
		return &ClientErrors.ServerBadChunks{} // Blame the server
	}

	// Creates a privte socket connection between the server to upload the file to the server
	uploadSocket, err := Helper.CreatePrivateSocket()
	if err != nil {
		return err
	}

	go uploadFile(int64(fileSize), chunksSize, filename, *uploadSocket)

	return nil
}

// Handles download command
func HandleDownloadFile(command_arguments []string, socket *net.Conn) error {
	if len(command_arguments) < minimumdownloadArguments {
		return &ClientErrors.InvalidArgumentCountError{Arguments: uint8(len(command_arguments)), Expected: uint8(operationArguments)}
	}
	var filename string
	var clientpath string
	if Helper.IsQuoted(command_arguments, Helper.OneClosedPath) { // Check if the first command argument is enclosed within a quotation (') marks
		filename = Helper.FindPath(command_arguments, Helper.FirstNameParameter, Helper.OneClosedPath) // Save the first path (filename to upload)
		filename = filename[Helper.SkipEnclose : len(filename)-Helper.SkipEnclose]                     // Remove enclouse chars
		if Helper.IsQuoted(command_arguments, Helper.TwoCloudPaths) {                                  // Check if the second command argument is enclosed within a quotation (') marks
			clientpath = Helper.FindPath(command_arguments, Helper.SecondNameParameter, Helper.TwoCloudPaths) // Save the second path (path in cloud storage to save)
		} else { // If first path is quoted but the second doesn't
			clientpath = Helper.ReturnNonQuotedSecondPath(command_arguments)
		}
	} else { // If command arguments are not enclosed within a quotation (') marks
		// relay on argument indexes
		filename = command_arguments[oldFileName]
		if len(command_arguments) >= localPathIndex { // If local path has been specified
			clientpath = command_arguments[newFileName]
		}
	}
	data, err := Helper.ConvertStringToBytes(filename)
	if err != nil {
		return err
	}

	respone, err := Requests.SendRequest(Requests.DownloadFileRequest, data, socket) // Sends download file request
	if err != nil {
		return err
	}

	chunksSize, err := Helper.ConvertResponeToChunks(respone) // Convert respone to chunks size
	if err != nil {                                           // If chunks size was returned from the server in a wrong type
		return &ClientErrors.ServerBadChunks{} // Blame the server
	}

	// Creates a privte socket connection between the server to download the file from the server
	downloadSocket, err := Helper.CreatePrivateSocket()
	if err != nil {
		return err
	}

	go downloadFile(filepath.Join(clientpath, filepath.Base(filename)), chunksSize, *downloadSocket)

	return nil
}
