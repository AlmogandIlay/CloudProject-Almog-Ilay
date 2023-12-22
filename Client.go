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

	//filename := getFileName()
	filename := "C:/Users/אילאי/OneDrive/מסמכים/checkfiles/200MB.bin"
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error is: ", err)
	}

	defer f.Close()

	serverAddr := "46.116.205.123:12345"
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
	fileSize := fileInfo.Size()
	println("Size of the file", fileSize)
	_, err = conn.Write([]byte(strconv.Itoa(int(fileSize))))
	if err != nil {
		fmt.Println(err.Error())
	}
	chunk := make([]byte, chunkSize)

	start := time.Now()

	var totalBytesRead int64
	var totalReadFlag int64
	var precentage int64
	for {
		bytesRead, err := f.Read(chunk)
		if err == io.EOF || bytesRead == 0 {
			break
		}
		if err != nil {
			println(err.Error())
		}

		_, err = conn.Write(chunk[:bytesRead])
		if err != nil {
			fmt.Println("Error sending file chunk:", err)
			os.Exit(1)
		}

		totalReadFlag += int64(bytesRead)
		totalBytesRead += int64(bytesRead)

		if totalReadFlag >= 1_000_000 {
			totalReadFlag = 0
			precentage = (totalBytesRead * 100) / fileSize

			printer := func(length int, char string) string {
				var s string
				for i := 0; i < int(precentage/2); i++ {
					s += "-"
				}
				return s
			}

			fmt.Printf("\r%v    Download progress: %v%% - %s", time.Since(start), precentage, printer(int(precentage), "-"))
		}
	}
	fmt.Print("\n")
	elapsed := time.Since(start)
	fmt.Printf("Binomial took %s\n", elapsed)
	fmt.Println("Success!")
}
