package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
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
		fmt.Println("Error is:", err)
	}

	defer f.Close()

	serverAddr := "46.116.199.220:12345"
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to the server!")
	fileInfo, err := os.Stat(filename)
	println("Size of the file", fileInfo.Size())
	println("length of the string converstion ", len(strconv.FormatInt(fileInfo.Size(), 10)))
	_, err = conn.Write([]byte(strconv.Itoa(int(fileInfo.Size()))))
	start := time.Now()

	buf := make([]byte, 1024*1024)
	var s []byte

	for {
		bytesRead, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		print(string(buf[:bytesRead]))
		s = append(s, buf[:bytesRead]...)
	}

	_, err = conn.Write([]byte(s))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	elapsed := time.Since(start)

	fmt.Print("\n\ntime:", elapsed, "\n")
	fmt.Println("Success!")
}
