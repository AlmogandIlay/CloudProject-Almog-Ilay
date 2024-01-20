package Menu

import (
	"fmt"
)

type CLI struct {
	prompt string
}

func NewCLI() *CLI {
	return &CLI{prompt: ">> "}
}

// Prints the program startup intro
func (cli *CLI) PrintStartup() {
	fmt.Println("CloudDrive v1.0 Command Line Interface")
	fmt.Println("Type \"help\" for available commands.")
}
