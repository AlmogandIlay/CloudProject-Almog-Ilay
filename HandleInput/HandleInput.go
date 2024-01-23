package Handleinput

import (
	"bufio"
	"client/Authentication"
	"fmt"
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

/*
Scan user's input and convert it to text
*/
func (inputBuffer UserInput) readInput() string {
	inputBuffer.Scanner.Scan()
	command := inputBuffer.Scanner.Text()

	return command
}

/*
Gets user input and handles its command request.
*/
func (inputBuffer UserInput) Handleinput(socket net.Conn) string {
	var err error
	command := strings.Fields(strings.ToLower(inputBuffer.readInput()))
	if len(command) > 0 { // If command is not empty
		command_prefix := command[prefix_index]

		switch command_prefix {

		case "signup":
			err = Authentication.HandleSignup(command[command_arguments:])
			if err != nil {
				fmt.Println(err.Error())
			}

		case "help":

		case "cd":

		}

		return ""
	}
	return ""
}
