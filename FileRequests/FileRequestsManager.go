package FileRequestsManager

import (
	"client/Helper"
	"client/Requests"
	"fmt"
	"net"
	"strings"
)

const (
	minimum_arguments = 1
	path_argument     = 0
	path_index        = 1
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
