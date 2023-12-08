package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	serverAddr := "192.168.50.191:12345"

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to the server!")

	fmt.Println("Welcome to our Cloud\nWhat would you like to do:\n1. LOGIN\n2. SIGN IN\nYour choice: ")

}
