package FileSystem

import (
	filetransmission "CloudDrive/FileTransmission"
	helper "CloudDrive/Helper"
	"net"
	"os"
	"path/filepath"
)

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

	return filetransmission.GetChunkSize(file.Size), nil
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

	err = IsContentInDirectory(dir.Name, dir.Path)
	if err == nil { // If dir exists
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
	return filetransmission.GetChunkSize(fileSize), fileSize, nil
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

	err = filetransmission.ReceiveFile(*uploadSocket, file.Path, file.Name, int(file.Size))
	if err != nil { // If upload process has failed
		err = sendResponseInfo(uploadSocket, buildError(err.Error())) // Send error respone
		if err != nil {
			return // Exit upload process
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

	_ = os.Mkdir(fullPath, os.ModePerm) // sets permissions for the directory

	err = filetransmission.ReceiveFolder(uploadSocket, fullPath)
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

	err = filetransmission.SendFile(downloadSocket, fileSize, filepath) // Send file to client
	if err != nil {
		return // Exit download process
	}
}
