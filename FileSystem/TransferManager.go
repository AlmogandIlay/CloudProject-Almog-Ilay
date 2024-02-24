package FileSystem

import (
	filetransmission "CloudDrive/FileTransmission"
	helper "CloudDrive/Helper"
	"net"
	"os"
	"path/filepath"
)

// Create private socket and ready to accept the client connection
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

// Uploading file proccess
func uploadAbsFile(file *File, uploadListener *net.Listener) {
	// Create private socket with the client for the upload file
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
			return // Exit upload
		}
	}
}
