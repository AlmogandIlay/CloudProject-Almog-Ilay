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
	fmt.Println("Json data is ", json)
	_, err = conn.Write([]byte(json))
	if err != nil {
		fmt.Println("Error caught: ", err.Error())
	}
	conn.Close()
	fmt.Println("Client 1 Disconnected")

	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to the server!")
	json = "{\"Type\": 102, \"Data\": {\"username\":\"Almog\",\"password\":\"SecretPas123\",\"email\":\"almog.eisen@gmail.com\"} }"
	fmt.Println("Json data is ", json)
	_, err = conn.Write([]byte(json))
	if err != nil {
		fmt.Println("Error caught: ", err.Error())
	}
	conn.Close()
	fmt.Println("Client 2 disconnected")

	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to the server!")
	json = "{\"Type\": 101, \"Data\": {\"username\":\"Almog\",\"password\":\"SecretPas123\",\"email\":\"almog.eisen@gmail.com\"} }"
	fmt.Println("Json data is ", json)
	_, err = conn.Write([]byte(json))
	if err != nil {
		fmt.Println("Error caught: ", err.Error())
	}
	conn.Close()
	fmt.Println("Client 3 disconnected")

	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to the server!")
	json = "{\"Type\": wrongType, \"Data\": {\"username\":\"Almog\",\"password\":\"SecretPas123\",\"email\":\"almog.eisen@gmail.com\"} }"
	fmt.Println("Json data is ", json)
	_, err = conn.Write([]byte(json))
	if err != nil {
		fmt.Println("Error caught: ", err.Error())
	}
	conn.Close()
	fmt.Println("Client 4 disconnected")

	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to the server!")
	json = "{\"Type\": 102, \"Data\": {\"username\":\"Almog\",\"password\":\"SecretPas123\",\"email\":\"\"} }"
	fmt.Println("Json data is ", json)
	_, err = conn.Write([]byte(json))
	if err != nil {
		fmt.Println("Error caught: ", err.Error())
	}
	conn.Close()
	fmt.Println("Client 5 disconnected")

	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to the server!")
	json = "{\"Type\": 102, \"Data\": {\"username\":\"Almog\",\"password\":\"SecretPas123\",\"email\":} }"
	fmt.Println("Json data is ", json)
	_, err = conn.Write([]byte(json))
	if err != nil {
		fmt.Println("Error caught: ", err.Error())
	}
	conn.Close()
	fmt.Println("Client 6 disconnected")

	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to the server!")
	json = "{\"Type\": 102, \"Data\": {\"username\":\"Almog\",\"password\":\"SecretPas123\",} }"
	fmt.Println("Json data is ", json)
	_, err = conn.Write([]byte(json))
	if err != nil {
		fmt.Println("Error caught: ", err.Error())
	}
	conn.Close()
	fmt.Println("Client 7 disconnected")

	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to the server!")
	json = "{\"Type\": 102, \"Data\": {\"username\":\"Almog\",\"password\":\"SecretPas123\"} }"
	fmt.Println("Json data is ", json)
	_, err = conn.Write([]byte(json))
	if err != nil {
		fmt.Println("Error caught: ", err.Error())
	}
	conn.Close()
	fmt.Println("Client 8 disconnected")

	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to the server!")
	json = "{\"Type\": 102, \"Data\": }"
	fmt.Println("Json data is ", json)
	_, err = conn.Write([]byte(json))
	if err != nil {
		fmt.Println("Error caught: ", err.Error())
	}
	conn.Close()
	fmt.Println("Client 9 disconnected")

	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to the server!")
	json = "{\"Type\": 102, \"Data\": {} }"
	fmt.Println("Json data is ", json)
	_, err = conn.Write([]byte(json))
	if err != nil {
		fmt.Println("Error caught: ", err.Error())
	}
	conn.Close()
	fmt.Println("Client 10 disconnected")

	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to the server!")
	json = "{\"Type\": 102, }"
	fmt.Println("Json data is ", json)
	_, err = conn.Write([]byte(json))
	if err != nil {
		fmt.Println("Error caught: ", err.Error())
	}
	conn.Close()
	fmt.Println("Client 11 disconnected")

	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to the server!")
	json = "{\"Type\": 102 }"
	fmt.Println("Json data is ", json)
	_, err = conn.Write([]byte(json))
	if err != nil {
		fmt.Println("Error caught: ", err.Error())
	}
	conn.Close()
	fmt.Println("Client 12 disconnected")

}
