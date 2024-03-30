package main

import (
	Menu "client/Menu"
	"fmt"
	"log"
)

func main() {
	cli, err := Menu.NewCLI()
	if err != nil { // If server connection fails
		fmt.Printf("\n") // Divide connection error from other texts
		log.Fatal(err)
	}

	cli.PrintStartup()
	cli.Loop()
}
