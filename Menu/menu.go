package Menu

import (
	"errors"
	"fmt"
	"net"
)

const (
	ip_addr = "192.168.50.191:12345"
)

type CLI struct {
	prompt string
	socket net.Conn
}

func NewCLI() (*CLI, error) {
	var err error
	var cli CLI

	cli.prompt = ">>"
	cli.socket, err = net.Dial("tcp", ip_addr)
	if err != nil {

		return nil, errors.New(fmt.Sprintf("There has been an error connecting to the server.\nPlease check your connection and try again.\nIf it doesn't work contact the developers and send them this error message:\n%s", err.Error()))
	}

	return &cli, nil

}

// Prints the program startup intro
func (cli *CLI) PrintStartup() {
	fmt.Println("CloudDrive v1.0 Command Line Interface")
	fmt.Println("Type \"help\" for available commands.")
}
