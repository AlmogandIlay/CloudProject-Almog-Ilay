package main

import "client/menu"

func main() {
	cli := menu.NewCLI()
	cli.PrintStartup()
}
