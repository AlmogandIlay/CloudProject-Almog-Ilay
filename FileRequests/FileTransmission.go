package FileRequestsManager

import (
	"bufio"
	"client/ClientErrors"
	"client/Helper"
	"client/Requests"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/galsondor/go-ascii"
)

const (
	empty    = 0
	kilobyte = 1_000_000

	firstIndex    = 0
	onlyCharacter = 1
)

//var mu sync.Mutex // Lock the file writing to make sure only one goroutine can write over the file

// Uploads file to the cloud server
func uploadFile(fileSize int64, chunksSize int, filename string, socket net.Conn) {
	file, err := os.Open(strings.Replace(filename, "'", "", Helper.RemoveAll)) // Open file
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	chunk := make([]byte, chunksSize) // Save buffer of chunks

	var validUpload = true
	var totalBytesRead int64
	var totalReadFlag int64 // Flag for client view to automatically update upload percentage and bar progress
	var precentage int64

	for {
		bytesRead, err := file.Read(chunk)
		if err == io.EOF { // If finish reading file succesfully
			break
		}
		if err != nil { // If error occurred while reading the file
			validUpload = false
			fmt.Println(err.Error())
			err = &ClientErrors.ReadFileInfoError{}
			fmt.Println(err.Error())
		}
		if bytesRead == empty { // If finish reading file succesfully
			break
		}

		_, err = socket.Write(chunk[:bytesRead]) // Sending chunk to server
		if err != nil {                          // If sending error occured
			validUpload = false
			err = &ClientErrors.SendDataError{}
			fmt.Println(err.Error())
		}

		// Add bytes read for the chunk to the flag percentage and to the total read bytes
		totalReadFlag += int64(bytesRead)
		totalBytesRead += int64(bytesRead)

		if totalReadFlag >= kilobyte { // For every 1 Kilobyte update the progess and perecntage bar
			totalReadFlag = 0
			precentage = (totalBytesRead * 100) / fileSize // Calculates total read bytes compared to the total file size in percentages
			printer := func(_ int, _ string) string {
				var bar string
				for i := 0; i < int(precentage/2); i++ {
					bar += "-"
				}
				return bar
			}
			fmt.Printf("\033[F\033[K")
			fmt.Printf("Upload Progress: %v%% - %s", precentage, printer(int(precentage), "-"))
			fmt.Println()

		}
	}
	if validUpload {
		fmt.Printf("File %s has been uploaded successfully\n", filename)
	}
}

// Upload directory to cloud server
func uploadDirectory(dirpath string, socket net.Conn) {
	err := filepath.WalkDir(dirpath, func(contentPath string, contentInfo fs.DirEntry, err error) error { // Walk through all the contents in the given dir path
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(dirpath, contentPath) // Convert the contentPath absolute to relative from the given path to upload
		if err != nil {
			return &ClientErrors.ConvertToRelative{}
		}

		if relativePath != "." { // If path is not the base (already exists) path
			dirData, err := Helper.ConvertStringToBytes(relativePath) // Convert new dir path to bytes
			if err != nil {
				return err
			}

			if !contentInfo.IsDir() { // If content is file

				fileInfo, err := contentInfo.Info() // Get file's info
				if err != nil {
					return &ClientErrors.ReadFileInfoError{Filename: filepath.Base(relativePath)}
				}

				// Initializes file struct
				file := newContent(filepath.Base(relativePath), filepath.Dir(relativePath), uint32(fileInfo.Size()))
				// Convert file struct to json bytes
				file_data, err := json.Marshal(file)
				if err != nil {
					return &ClientErrors.JsonEncodeError{Err: err}
				}

				// Sends Upload File reques
				respone, err := Requests.SendRequest(Requests.UploadFileRequest, file_data, &socket)
				if err != nil { // If upload file request was rejected
					return err
				}

				chunksSize, err := Helper.ConvertResponeToChunks(respone) // Convert respone to chunks size
				if err != nil {                                           // If chunks size was returned from the server in a wrong type
					return &ClientErrors.ServerBadChunks{} // Blame the server
				}

				uploadFile(fileInfo.Size(), chunksSize, contentPath, socket) // Uploads the file

			} else { // If content is directory
				// Sends request to make a new directory
				respone, err := Requests.SendRequestInfo(Requests.BuildRequestInfo(Requests.CreateFolderRequest, dirData), socket)
				if err != nil {
					return err
				}

				if respone.Type == Requests.ErrorRespone { // If respone is error
					return fmt.Errorf(respone.Respone)
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
	} else { // If upload was valid
		_, err = Requests.SendRequestInfo(Requests.BuildRequestInfo(Requests.StopUpload, nil), socket) // Send stop upload request to server
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Upload directory has finished")
	}
}

// Write the file content on a seprated goroutine to not waste time and resources for the main thread that recives the file
// func writeFile(writer **bufio.Writer, chunkBytes []byte) {
// 	fmt.Println("Locking... (Currently Unlocked)")
// 	mu.Lock()         // Acquire lock before writing the file, promising that only one goroutine can write over the file
// 	defer mu.Unlock() // Ensure the lock is released after writing the file
// 	fmt.Println("It's locked")
// 	_, err := (*writer).Write(chunkBytes)
// 	if err != nil {
// 		fmt.Println("Error writing data on the provided file path.\nPlease contact the developers")
// 		return
// 	}
// 	fmt.Println("Unlocking... (Currently locked)")
// }

// Download a file from the cloud server
func downloadFile(path string, chunksSize int, socket net.Conn) {
	file, err := os.Create(path) // Creates the file in the given/default path
	if err != nil {
		fmt.Println("Couldn't create the file in the provided path.\nPlease provide a different path.")
		return
	}
	file.Close()

	file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644) // Open file for writing
	if err != nil {
		fmt.Println("Couldn't open the file for writing data.\nPlease make sure you have provided a path with permissions for write.")
		return
	}
	defer file.Close()

	// Create a buffered writier for efficient writes
	writer := bufio.NewWriter(file)

	for {
		chunkBytes, err := Helper.ReciveChunkData(&socket, chunksSize)
		// If server indicated that the end of the file has reached OR the client hasn't recived any new chunks for over the configured timeout, finish reading file sucessfully
		if netErr, ok := err.(*net.OpError); ok && netErr.Timeout() || (chunkBytes[firstIndex] == ascii.ETX && len(chunkBytes) == onlyCharacter) {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		_, err = (*writer).Write(chunkBytes)
		if err != nil {
			fmt.Println("Error writing data on the provided file path.\nPlease contact the developers")
			return
		}
		//go writeFile(&writer, chunkBytes) // Write the file over goroutine to not interept the connection
	}
	err = writer.Flush() // Flush any remaining data in the buffer to the file
	if err != nil {
		fmt.Println("Error flushing data to the file.\nPlease contact the developers")
		return
	}
	fmt.Printf("File %s has been downloaded successfully\n", path)
}
