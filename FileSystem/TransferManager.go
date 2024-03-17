package FileSystem

import (
	"CloudDrive/FileTransmission"
	helper "CloudDrive/Helper"
	"CloudDrive/Server/RequestHandlers/Requests"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

const (
	okayRespone   string = "Okay"
	chunksRespone string = "ChunksSize:"
)

type clientResponeInfo struct {
	Type    int    `json:"Type"`
	Respone string `json:"Data"`
}

// Create private socket for file recieve and ready to accept the client connection
func createPrivateSocket(uploadListener net.Listener) (*net.Conn, error) {
	conn, err := uploadListener.Accept()
	if err != nil {
		return nil, &CreatePrivateSocketError{}
	}
	return &conn, nil
}

// Upload a file to the Cloud
func (user *LoggedUser) UploadFile(file *Content, uploadListener *net.Listener) (uint, error) {
	if file.Path == "" { // if path wasn't decleared
		file.Path = user.GetPath()
	}
	err := user.ValidateContent(*file) // validate file content
	if err != nil {
		return emptyChunks, err
	}

	err = user.validNewContentSize(file.Size, file.Name) // validate new file size. Checks if the current storage space can handle the file
	if err != nil {
		return emptyChunks, err
	}

	if !filepath.IsAbs(file.Path) { // Convert file's path to absolute if it doesn't
		file.Path = helper.ConvertToAbsolute(user.GetPath(), file.Path)
	}

	err = IsContentInDirectory(file.Name, file.Path)
	if err == nil { // If file exists
		return emptyChunks, &FileExistError{file.Name, file.Path}
	}

	go uploadAbsFile(file, uploadListener) // Start receiving file from client

	return FileTransmission.GetChunkSize(file.Size), nil
}

func (user *LoggedUser) DownloadDirectory(dirPath string, downloadListener *net.Listener) error {

	if !helper.IsAbs(dirPath) { // if dir path is relative
		dirPath = helper.ConvertToAbsolute(user.GetPath(), dirPath) // Convert the given path to absolute server-side path
	}
	dirPath = helper.GetServerStoragePath(user.UserID, dirPath)              // Convert path to server-side path in case it was an absolute path
	err := IsContentInDirectory(helper.Base(dirPath), filepath.Dir(dirPath)) // Check if dir does exist in the directory
	if err != nil {
		return err
	}

	socket, err := createPrivateSocket(*downloadListener)
	if err != nil {
		return err
	}
	go downloadAbsDirectory(dirPath, socket) // Start sending file to client

	return nil
}

func (user *LoggedUser) UploadDirectory(dir *Content, uploadListener *net.Listener) error {
	if dir.Path == "" { // if path wasn't decleared
		dir.Path = user.GetPath()
	}
	err := user.ValidateContent(*dir) // validate dir info
	if err != nil {
		return err
	}

	err = user.validNewContentSize(dir.Size, dir.Name) // validate new dir size. Checks if the current storage space can handle the dir
	if err != nil {
		return err
	}

	if !filepath.IsAbs(dir.Path) { // Convert dir's path to absolute if it doesn't
		dir.Path = helper.ConvertToAbsolute(user.GetPath(), dir.Path)
	}

	err = IsContentInDirectory(dir.Name, dir.Path) // Checks if dir doesn't exist already
	if err == nil {                                // If dir exists
		return &ContentExistError{dir.Name, dir.Path}
	}

	go uploadAbsDirectory(dir, uploadListener) // Start receiving directory from client
	return nil
}

func (user *LoggedUser) DownloadFile(filePath string, downloadListener *net.Listener) (uint, uint32, error) {
	if !helper.IsAbs(filePath) { // if file path is relative
		filePath = helper.ConvertToAbsolute(user.GetPath(), filePath) // Convert filepath to absolute server-side path
	}
	filePath = helper.GetServerStoragePath(user.UserID, filePath)              // Convert path to server-side path in case it was an absolute path
	err := IsContentInDirectory(helper.Base(filePath), filepath.Dir(filePath)) // Check if file does exist in the directory
	if err != nil {
		return emptyChunks, emptySize, err
	}

	go downloadAbsFile(filePath, downloadListener) // Start sending file to client

	fileSize, err := getFileSize(filePath)
	if err != nil {
		return emptyChunks, emptySize, err
	}
	return FileTransmission.GetChunkSize(fileSize), fileSize, nil
}

// Uploading file proccess
func uploadAbsFile(file *Content, uploadListener *net.Listener) {
	// Creates a private socket with the client for the upload file
	uploadSocket, err := createPrivateSocket(*uploadListener)
	if err != nil {
		return // Exit upload proccess
	}

	fullPath := file.Path + "\\" + file.Name // Saves the full path for the file to be created
	dirFile, _ := os.Create(fullPath)        // Creates the file
	dirFile.Close()

	err = FileTransmission.ReceiveFile(uploadSocket, file.Path, file.Name, int(file.Size))
	if err != nil { // If upload process has failed
		err = sendResponseInfo(uploadSocket, buildError(err.Error())) // Send error respone
		if err != nil {
			return // Exit upload process
		}
	}
}

// Implement createFolder for reciving folder implemention.
// Input:
// info Requests.RequestInfo - Client request's info
// baseFolderPath - Base absolute path of server side
func createFolder(info Requests.RequestInfo, baseFolderPath string) clientResponeInfo {
	// Convert request json bytes to dir path variable
	rawData := Requests.ParseDataToString(info.RequestData)
	relativeFolderPath := helper.ConvertRawJsonToData(rawData)
	absFolderPath := filepath.Join(baseFolderPath, relativeFolderPath) // Appened the base path and the relative path to make a full absolute path

	err := IsContentInDirectory(helper.Base(absFolderPath), filepath.Dir(absFolderPath)) // Check if the folder is already exists
	if err == nil {                                                                      // If directory is already exist
		return clientResponeInfo{Type: errorRespone, Respone: (&FolderExistError{Name: helper.Base(absFolderPath), Path: filepath.Dir(absFolderPath)}).Error()} // If folder exists
	}
	os.Mkdir(absFolderPath, os.ModePerm) // Creates a folder
	return clientResponeInfo{Type: validRespone, Respone: okayRespone}
}

// Implement createFile for reciving folder implemention.
// Input:
// Requests.RequestInfo - Client request's info
// baseFolderPath - Base absolute path of server side
// Output:
// Respone (ResponeInfo)
// Absolute filepath (if file's valid)
// File's size (if file's valid)
func createFile(info Requests.RequestInfo, baseFolderPath string) (clientResponeInfo, string, int) {
	file, err := ParseDataToContent(info.RequestData) // Parse json bytes to file struct
	if err != nil {
		return clientResponeInfo{Type: errorRespone, Respone: err.Error()}, "", nonSize
	}
	err = validContentName(file.Name) // Valids file name
	if err != nil {
		return clientResponeInfo{Type: errorRespone, Respone: err.Error()}, "", nonSize
	}
	absFilePath := filepath.Join(baseFolderPath, file.Path, file.Name) // Appened the base path and the relative path to make a full absolute path

	err = IsContentInDirectory(helper.Base(absFilePath), filepath.Dir(absFilePath)) // Check if the file is already exists
	if err == nil {                                                                 // If file exists
		// Returns FileExist error
		return clientResponeInfo{Type: errorRespone, Respone: (&FileExistError{Name: helper.Base(absFilePath), Path: filepath.Dir(absFilePath)}).Error()}, "", nonSize
	}

	chunksSize := FileTransmission.GetChunkSize(file.Size) // Saves the chunks size for the file

	dirFile, _ := os.Create(absFilePath) // Creates the file
	dirFile.Close()

	return clientResponeInfo{Type: validRespone, Respone: chunksRespone + strconv.FormatUint(uint64(chunksSize), 10)}, absFilePath, int(file.Size)
}

// Recieves folder from client implemention
func receiveFolder(conn *net.Conn, absDirPath string) error {
	for {
		// File variables:
		var absFilePath string // Absolute file path
		var fileSize int       // File size

		request_Info, err := Requests.ReciveRequestInfo(conn, true) // Recieves request info with timeout flag on, indicating whether to upload file or folder
		if err != nil {
			switch err := err.(type) { // Checking error type
			case *net.OpError:
				if err.Timeout() { // If error is reciving timeout
					return &UploadTimeOut{}
				}
			}
		}
		var responeInfo clientResponeInfo
		switch request_Info.Type {
		case Requests.CreateFolderRequest:
			responeInfo = createFolder(request_Info, absDirPath) // Creates the directory with a given absolute base path, returns the respone info indicating success

		case Requests.UploadFileRequest:
			responeInfo, absFilePath, fileSize = createFile(request_Info, absDirPath) // Creates the file with a given absolute base path, returns the respone info including the file's chunk size
			fmt.Println(responeInfo.Respone)

		case Requests.StopUpload: // If client requested to stop uploading
			return nil
		}
		message, err := json.Marshal(responeInfo) // encode respone info to json bytes
		if err != nil {
			return err
		}
		helper.SendData(conn, message)                                                           // Sends respone Info bytes
		if request_Info.Type == Requests.UploadFileRequest && responeInfo.Type != errorRespone { // If client requested to upload a file and its respone is valid
			FileTransmission.ReceiveFile(conn, filepath.Dir(absFilePath), helper.Base(absFilePath), fileSize) // Start reciving file proccess
		}
	}
}

// Uploading directory proccess
func uploadAbsDirectory(dir *Content, uploadListener *net.Listener) {
	// Creates a private socket with the client for the upload directory
	uploadSocket, err := createPrivateSocket(*uploadListener)
	if err != nil {
		return // Exit upload proccess
	}

	fullPath := dir.Path + "\\" + dir.Name // Saves the full path for the directory to be created

	_ = os.Mkdir(fullPath, os.ModePerm) // Creates the base directory with set permissions for the directory on the cloud

	err = receiveFolder(uploadSocket, fullPath)
	if err != nil { // If upload process has failed
		sendResponseInfo(uploadSocket, buildError(err.Error())) // Send error respone
	}
}

func downloadAbsDirectory(directoryPath string, downloadSocket *net.Conn) {

	err := filepath.WalkDir(directoryPath, func(contentpath string, contentInfo fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relativePath, err := filepath.Rel(directoryPath, contentpath)
		if err != nil {
			return err // Todo: create custom error
		}

		if relativePath != "." {
			if !contentInfo.IsDir() {
				fileInfo, err := contentInfo.Info()

				if err != nil {
					return err // Todo: create custom error
				}

				file := newContent(helper.Base(relativePath), filepath.Dir(relativePath), uint32(fileInfo.Size()))

				fileData, err := json.Marshal(file)

				if err != nil {
					return &MarshalError{}
				}

				responeInfo := buildRespone(fileData)

				err = sendResponseInfo(downloadSocket, responeInfo)

				if err != nil {
					return err
				}

				chunkSize := FileTransmission.GetChunkSize(uint32(file.Size))

				responeInfo = buildRespone([]byte(chunksRespone + strconv.FormatUint(uint64(chunkSize), 10)))

				err = sendResponseInfo(downloadSocket, responeInfo)

				err = FileTransmission.SendFile(downloadSocket, uint64(file.Size), relativePath) // Send file to client
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Downloading file process
func downloadAbsFile(filepath string, downloadListener *net.Listener) {
	fileInfo, err := os.Stat(filepath)
	if err != nil { // Server-side error (very unlikely to happen)
		return
	}
	fileSize := uint64(fileInfo.Size())
	// Creates a private socket with the client for the download file
	downloadSocket, err := createPrivateSocket(*downloadListener)
	if err != nil {
		return // Exit download proccess
	}

	err = FileTransmission.SendFile(downloadSocket, fileSize, filepath) // Send file to client
	if err != nil {
		return // Exit download process
	}
}
