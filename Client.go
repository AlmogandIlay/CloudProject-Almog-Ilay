package main

import (
	Menu "client/Menu"
)

func main() {
	cli, err := Menu.NewCLI()
	if err != nil { // If server connection fails
		panic(err)
	}

	cli.PrintStartup()
	cli.Loop()
}
