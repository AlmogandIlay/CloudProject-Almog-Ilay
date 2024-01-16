package Handleinput

import (
	"bufio"
	"os"
	"strings"
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
	command := strings.ToLower(inputBuffer.convertToText()) // Recieves user command and saves it as case insensitive

	switch command {
	case "help":

	case "cd":

	}

	return ""
}
