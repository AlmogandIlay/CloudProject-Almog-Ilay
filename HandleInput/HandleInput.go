package Handleinput

import (
	"bufio"
	"client/Authentication"
	FileRequestsManager "client/FileRequests"
	"net"
	"os"
	"strings"
)

const (
	prefix_index      = 0
	command_arguments = 1
)

type UserInput struct {
	Scanner *bufio.Scanner
}

func NewUserInput() *UserInput {
	return &UserInput{Scanner: bufio.NewScanner(os.Stdin)}
}

// Scan user's input and convert it to text
func (inputBuffer UserInput) readInput() string {
	inputBuffer.Scanner.Scan()
	command := inputBuffer.Scanner.Text()

	return command
}

func helpScreen() string {
	help_command :=
		`
SIGNUP		Create an account in CloudDrive service.
SIGNIN		Sign in to an existing CloudDrive account.
CD		Displays/Changes the current working directory.
NEWFILE		Creates a new file.
NEWDIR		Creates a new directory.
RMFILE		Removes a file.
RMFOLDER		Removes a folder.	
RENAME		Renames a folder or a directory.
MOVE		Moves a file/folder to a different location.
LS		List all the current files in the current or given path.
`
	return help_command
}

//Gets user input and handles its command request.

func (inputBuffer UserInput) Handleinput(socket net.Conn) string {
	var err error
	command := strings.Fields(inputBuffer.readInput())
	if len(command) > 0 { // If command is not empty
		command_prefix := strings.ToLower(command[prefix_index])

		switch command_prefix {

		case "help":
			return helpScreen()

		case "signup":
			err = Authentication.HandleSignup(command[command_arguments:], socket)
			if err != nil {
				return err.Error()
			}

			FileRequestsManager.InitializeCurrentPath()
			return "Successfully signed up!\n"

		case "signin":
			err = Authentication.HandleSignIn(command[command_arguments:], socket)
			if err != nil {
				return err.Error()
			}

			FileRequestsManager.InitializeCurrentPath()
			return "Successfully signed in!\n"

		case "cd":
			err = FileRequestsManager.HandleChangeDirectory(command[command_arguments:], socket)
			if err != nil {
				return err.Error()
			}
			return ""

		case "newfile":
			err = FileRequestsManager.HandleCreate(command[command_arguments:], "createFile", socket)
			if err != nil {
				return err.Error()
			}
			return "File created successfully!\n"

		case "newdir":
			err = FileRequestsManager.HandleCreate(command[command_arguments:], "createFolder", socket)
			if err != nil {
				return err.Error()
			}
			return "Folder created successfully"

		case "rm file":
			err = FileRequestsManager.HandleRemove(command[command_arguments:], "removeFile", socket)
			if err != nil {
				return err.Error()
			}
			return "File deleted successfully!\n"

		case "rm folder":
			err = FileRequestsManager.HandleRemove(command[command_arguments:], "removeFolder", socket)
			if err != nil {
				return err.Error()
			}
			return "Folder deleted successfully!\n"
		case "rename":
			err = FileRequestsManager.HandleRename(command[command_arguments:], socket)
			if err != nil {
				return err.Error()
			}
			return "The content has been renamed!\n"
		case "move":
			err = FileRequestsManager.HandleMove(command[command_arguments:], socket)
			if err != nil {
				return err.Error()
			}
			return ""

		case "ls":
			dir, err := FileRequestsManager.HandleShow(command[command_arguments:], socket)
			if err != nil {
				return err.Error()
			}
			return dir

		default:
			return "Invalid command.\nPlease try a different command or \"help\"\n"

		}
	}
	return ""

}
