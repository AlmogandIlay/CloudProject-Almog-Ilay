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
func (user *LoggedUser) UploadFile(file *File, uploadListener *net.Listener) (uint, error) {
	if file.Path == "" { // if path wasn't decleared
		file.setPath(user.GetPath())
	}
	err := user.ValidateFile(*file) // validate file content
	if err != nil {
		return emptyChunks, err
	}

	err = user.validNewFileSize(file.Size) // validate new file size. Checks if the current storage space can handle the file
	if err != nil {
		return emptyChunks, err
	}

	if !filepath.IsAbs(file.Path) { // Convert file's path to absolute if it doesn't
		file.setPath(helper.ConvertToAbsolute(user.GetPath(), file.Path))
	}

	err = IsContentInDirectory(file.Name, file.Path)
	if err == nil { // If file exists
		return emptyChunks, &FileExistError{file.Name, file.Path}
	}

	go uploadAbsFile(file, uploadListener) // Start uploading file

	return filetransmission.GetChunkSize(file.Size), nil
}

func (user *LoggedUser) DownloadFile(filename string, downloadListener *net.Listener) error {
	if !helper.IsAbs(filename) { // if file path is relative
		filename = helper.ConvertToAbsolute(user.GetPath(), filename) // Convert filepath to absolute
	}
	realPath := helper.GetServerStoragePath(user.UserID, filename)
	err := IsContentInDirectory(helper.Base(realPath), filepath.Dir(realPath))
	if err != nil {
		return err
	}

	go downloadAbsFile(realPath, downloadListener) // Start downloading file

	return nil
}

// Uploading file proccess
func uploadAbsFile(file *File, uploadListener *net.Listener) {
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
