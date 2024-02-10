package FileRequestsManager

import (
	"client/Helper"
	"client/Requests"
	"fmt"
	"net"
)

const (
	minimum_arguments = 1
	path_argument     = 0
)

func HandleChangeDirectory(command_arguments []string, socket net.Conn) error {
	if len(command_arguments) < minimum_arguments {
		return fmt.Errorf("incorrect number of arguments.\nPlease try again")
	}

	data := Helper.ConvertStringToBytes(command_arguments[path_argument])
	err := Requests.SendRequest(Requests.ChangeDirectoryRequest, data, socket)

	return err
}
