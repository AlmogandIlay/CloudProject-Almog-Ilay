package Handleinput

import (
	"bufio"
	"os"
)

type UserInput struct {
	scanner *bufio.Scanner
}

func NewUserInput() *UserInput {
	return &UserInput{scanner: bufio.NewScanner(os.Stdin)}
}

// Scan input from client
func (input UserInput) GetInput() string {
	input.scanner.Scan()
	command := input.scanner.Text()

	return command

}

func Handleinput(command string)
