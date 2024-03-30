package main

import (
	Menu "client/Menu"
	"fmt"
	"log"
	"os"
)

func main() {
	cli, err := Menu.NewCLI()
	if err != nil { // If server connection fails
		log.Printf("\n%v\n\n", err)
		fmt.Println("Press 'Enter' to exit...")
		fmt.Scanln()
		os.Exit(1)
	}

	cli.PrintStartup()
	cli.Loop()
}
