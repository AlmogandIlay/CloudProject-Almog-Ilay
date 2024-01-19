package Handleinput

import (
	"bufio"
	"client/Authentication"
	"os"
	"strings"
)

const (
	prefix_index      = 0
	command_arguments = 1
)

type UserInput struct {
	scanner *bufio.Scanner
}

func NewUserInput() *UserInput {
	return &UserInput{scanner: bufio.NewScanner(os.Stdin)}
}

/*
Scans user's input and convert it to text
*/
func (inputBuffer UserInput) convertToText() string {
	inputBuffer.scanner.Scan()
	command := inputBuffer.scanner.Text()

	return command
}

/*
Gets user input and handles its command request.
*/
func (inputBuffer UserInput) Handleinput() string {
	command := strings.Fields(strings.ToLower(inputBuffer.convertToText()))
	command_prefix := command[prefix_index]

	switch command_prefix {

	case "signup":
		Authentication.HandleSignup(command[command_arguments:])

	case "help":

	case "cd":

	}

	return ""
}
