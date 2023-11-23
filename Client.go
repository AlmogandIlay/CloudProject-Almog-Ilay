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
		fmt.Println("Error is: ", err)
	}

	defer f.Close()
	buf := make([]byte, 1024*1024)
	s := make([]byte, 0, 1)

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
	for {
		_, err = f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("\n\n\n\n\n\n\n\n", err.Error(), "\n\n\n\n\n\n\n\n")
			break
		} else {
			fmt.Println(string(buf[:]))
		}
		s = append(s, buf...)
	}

	_, err = conn.Write([]byte(s))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	elapsed := time.Since(start)

	fmt.Print("time: ", elapsed, "\n")
	fmt.Println("Success!")
}
