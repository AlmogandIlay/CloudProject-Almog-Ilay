package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	chunkSize = 1024
)

func getFileName() string {
	fmt.Print("Enter file name: ")
	reader := bufio.NewReader(os.Stdin)
	filename, _ := reader.ReadString('\n')
	return strings.TrimSpace(filename)
}

func main() {

	filename := getFileName()
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error is: ", err)
	}

	defer f.Close()

	serverAddr := "localhost:12345"
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to the server!")

	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Println("Error getting file info:", err)
		os.Exit(1)
	}
	println("Size of the file", fileInfo.Size())
	_, err = conn.Write([]byte(strconv.Itoa(int(fileInfo.Size()))))

	chunk := make([]byte, chunkSize)
	for {
		bytesRead, err := f.Read(chunk)
		if err == io.EOF || bytesRead == 0 {
			break
		}
		if err != nil {
			println(err.Error())
		}
		println(string(chunk[:bytesRead]))
		_, err = conn.Write(chunk[:bytesRead])
		if err != nil {
			fmt.Println("Error sending file chunk:", err)
			os.Exit(1)
		}
	}

	fmt.Println("Success!")
}
