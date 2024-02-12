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

func HandleCreateFile(command_arguments []string, socket net.Conn) error {
	if len(command_arguments) != minimum_arguments {
		return fmt.Errorf("incorrect number of arguments.\nPlease try again")
	}

	data, err := Helper.ConvertStringToBytes(command_arguments[path_argument])
	if err != nil {
		return err
	}

	_, err = Requests.SendRequest(Requests.CreateFileRequest, data, socket)
	return err
}
