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

	fmt.Println("Connected to the server!")
	json := "{\"Type\": 101, \"Data\": {\"username\":\"almog\",\"password\":\"dfhdsfdfdffd\",\"email\":\"dfdsfjdsfj@sdfdfsjdsfjsdf\"} }"
	fmt.Println("Json data is", json)
	_, err = conn.Write([]byte(json))
	conn.Close()
	fmt.Println("Disconnected")
}
