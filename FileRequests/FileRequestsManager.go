package FileRequestsManager

import (
	"client/Helper"
	"client/Requests"
	"fmt"
	"net"
	"strings"
)

const (
	path_argument      = 0
	minimum_arguments  = 1
	rename_arguments   = 2
	oldFileName        = 0
	newFileName        = 1
	move_arguments     = 2
	show_abs_arguments = 1
	path_index         = 1
)

func convertResponeToPath(data string) string {
	parts := strings.Split(data, "CurrentDirectory:")
	return parts[path_index]
}

func HandleChangeDirectory(command_arguments []string, socket net.Conn) error {
	if len(command_arguments) < minimum_arguments {
		return fmt.Errorf("incorrect number of arguments.\nPlease try again")
	}

	data, err := Helper.ConvertStringToBytes(command_arguments[path_argument])
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

func HandleCreate(command_arguments []string, request string, socket net.Conn) error {
	if len(command_arguments) != minimum_arguments {
		return fmt.Errorf("incorrect number of arguments.\nPlease try again")
	}
	data, err := Helper.ConvertStringToBytes(command_arguments[path_argument])
	if err != nil {
		return err
	}

	switch request {
	case "createFile":
		return handleCreateFile(data, socket)
	case "createFolder":
		return handleCreateFolder(data, socket)
	}

	return fmt.Errorf("wrong create request")
}

func handleCreateFile(data []byte, socket net.Conn) error {
	_, err := Requests.SendRequest(Requests.CreateFileRequest, data, socket)
	return err
}

func handleCreateFolder(data []byte, socket net.Conn) error {
	_, err := Requests.SendRequest(Requests.CreateFolderRequest, data, socket)
	return err
}

func HandleRemove(command_arguments []string, request string, socket net.Conn) error {
	if len(command_arguments) != minimum_arguments {
		return fmt.Errorf("incorrect number of arguments.\nPlease try again")
	}
	data, err := Helper.ConvertStringToBytes(command_arguments[path_argument])
	if err != nil {
		return err
	}

	switch request {
	case "removeFile":
		return handleRemoveFile(data, socket)
	case "removeFolder":
		return handleRemoveFolder(data, socket)
	}

	return fmt.Errorf("wrong remove request")
}

func handleRemoveFile(data []byte, socket net.Conn) error {
	_, err := Requests.SendRequest(Requests.DeleteFileRequest, data, socket)
	return err
}

func handleRemoveFolder(data []byte, socket net.Conn) error {
	_, err := Requests.SendRequest(Requests.DeleteFolderRequest, data, socket)
	return err
}
func HandleRename(command_arguments []string, socket net.Conn) error {
	if len(command_arguments) != rename_arguments {
		return fmt.Errorf("incorrect number of arguments.\nPlease try again")
	}
	data, err := Helper.ConvertStringToBytes((command_arguments[oldFileName] + " " + command_arguments[newFileName]))
	if err != nil {
		return err
	}
	_, err = Requests.SendRequest(Requests.RenameRequest, data, socket)
	return err
}

func HandleMove(command_arguments []string, socket net.Conn) error {
	if len(command_arguments) != move_arguments {
		return fmt.Errorf("incorrect number of arguments.\nPlease try again")
	}
	data, err := Helper.ConvertStringToBytes(command_arguments[path_argument])
	if err != nil {
		return err
	}
	_, err = Requests.SendRequest(Requests.MoveRequest, data, socket)
	return err
}

func HandleShow(command_arguments []string, socket net.Conn) (string, error) {
	if len(command_arguments) != show_abs_arguments && len(command_arguments) != 0 { // check for amount of arguments
		return "", fmt.Errorf("incorrect number of arguments.\nPlease try again")
	}
	var data []byte
	var err error
	if len(command_arguments) == show_abs_arguments { // If specific path has been specified
		data, err = Helper.ConvertStringToBytes(command_arguments[path_argument])
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
