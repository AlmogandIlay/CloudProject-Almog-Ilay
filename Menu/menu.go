package Menu

import (
	FileRequestsManager "client/FileRequests"
	HandleInput "client/HandleInput"
	"fmt"
	"net"
)

const (
	ip_addr = "192.168.50.220:12345"
	prompt  = ">> "
)

type CLI struct {
	socket net.Conn
	prompt string
	input  *HandleInput.UserInput
}

func NewCLI() (*CLI, error) {

	sock, err := net.Dial("tcp", ip_addr)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("There has been an error connecting to the server.\nPlease check your connection and try again.\nIf it doesn't work contact the developers and send them this error message:\n%s", err.Error()))
	}
	return &CLI{socket: sock, prompt: prompt, input: HandleInput.NewUserInput()}, nil
}

func (cli *CLI) closeConnection() error {
	err := cli.socket.Close()
	if err != nil {
		return err
	}
	return nil
}

// Prints the program startup intro
func (cli *CLI) PrintStartup() {
	fmt.Println("CloudDrive v1.0 Command Line Interface!")
	fmt.Println("Type \"help\" for available commands.")
}

func (cli *CLI) printPrompt() {
	if FileRequestsManager.IsCurrentPathInitialized() {
		FileRequestsManager.PrintCurrentPath()
	}
	fmt.Print(cli.prompt)
}

func (cli *CLI) readInput() {
	cli.printPrompt()
	fmt.Println(cli.input.HandleInput(cli.socket))
}

func (cli *CLI) Loop() {
	defer cli.closeConnection()
	for {
		cli.readInput()
		if cli.input.Scanner.Bytes() == nil { // If unexpected input given
			break
		}
	}
}
