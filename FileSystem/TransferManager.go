package FileSystem

import (
	"CloudDrive/Filetransmission"
	helper "CloudDrive/Helper"
	"CloudDrive/Server/RequestHandlers/Requests"
	"encoding/json"
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

	return Filetransmission.GetChunkSize(file.Size), nil
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
	return Filetransmission.GetChunkSize(fileSize), fileSize, nil
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

	err = Filetransmission.ReceiveFile(*uploadSocket, file.Path, file.Name, int(file.Size))
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
func createFile(info Requests.RequestInfo, baseFolderPath string) clientResponeInfo {
	file, err := ParseDataToContent(info.RequestData) // Parse json bytes to file struct
	if err != nil {
		return clientResponeInfo{Type: errorRespone, Respone: err.Error()}
	}
	err = validContentName(file.Name) // Valids file name
	if err != nil {
		return clientResponeInfo{Type: errorRespone, Respone: err.Error()}
	}
	absFilePath := filepath.Join(baseFolderPath, file.Path, file.Name) // Appened the base path and the relative path to make a full absolute path

	err = IsContentInDirectory(helper.Base(absFilePath), filepath.Dir(absFilePath)) // Check if the file is already exists
	if err == nil {                                                                 // If file exists
		return clientResponeInfo{Type: errorRespone, Respone: (&FileExistError{Name: helper.Base(absFilePath), Path: filepath.Dir(absFilePath)}).Error()} // If file exists
	}

	chunksSize := Filetransmission.GetChunkSize(file.Size) // Saves the chunks size for the file

	dirFile, _ := os.Create(absFilePath) // Creates the file
	dirFile.Close()

	return clientResponeInfo{Type: validRespone, Respone: chunksRespone + strconv.FormatUint(uint64(chunksSize), 10)}
}

// Recieves folder from client implemention
func receiveFolder(conn *net.Conn, absDirPath string) error {
	for {
		request_Info, err := Requests.ReciveRequestInfo(conn) // Recieves request info indicating whether to upload file or folder
		if err != nil {
			return err
		}
		var responeInfo clientResponeInfo
		switch request_Info.Type {
		case Requests.CreateFolderRequest:
			responeInfo = createFolder(request_Info, absDirPath) // Creates the directory with a given absolute base path, returns the respone info indicating success

		case Requests.UploadFileRequest:
			responeInfo = createFile(request_Info, absDirPath) // Creates the file with a given absolute base path, returns the respone info including the file's chunk size
		}
		message, err := json.Marshal(responeInfo) // encode respone info to json bytes
		if err != nil {
			return err
		}
		helper.SendData(conn, message) // Sends respone Info bytes
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

	_ = os.Mkdir(fullPath, os.ModePerm) // Creates the base directory with set permissions for the directory

	err = receiveFolder(uploadSocket, fullPath)
	if err != nil { // If upload process has failed
		err = sendResponseInfo(uploadSocket, buildError(err.Error())) // Send error respone
		if err != nil {
			return // Exit upload process
		}
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

	err = Filetransmission.SendFile(downloadSocket, fileSize, filepath) // Send file to client
	if err != nil {
		return // Exit download process
	}
}
