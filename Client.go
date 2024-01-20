package main

import (
	connection "client/Connection"
	"client/Menu"
	"net"
)

var Socket *net.Conn

func main() {
	connection.Connect()

	cli := Menu.NewCLI()
	cli.PrintStartup()
}
