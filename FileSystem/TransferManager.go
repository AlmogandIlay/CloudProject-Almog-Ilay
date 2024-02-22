package FileSystem

import (
	filetransmission "CloudDrive/FileTransmission"
	helper "CloudDrive/Helper"
	"net"
	"os"
	"path/filepath"
)

// Upload a file to the Cloud
func (user *LoggedUser) UploadFile(file *File, conn *net.Conn) (uint, error) {
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

	go uploadAbsFile(file, conn) // Start uploading file

	return filetransmission.GetChunkSize(file.Size), nil
}

// Uploading file proccess
func uploadAbsFile(file *File, conn *net.Conn) {
	fullPath := file.Path + "\\" + file.Name // Saves the full path for the file to be created
	dirFile, _ := os.Create(fullPath)        // Creates the file
	dirFile.Close()

	err := filetransmission.ReceiveFile(*conn, file.Path, file.Name, int(file.Size))
	if err != nil { // If upload process has failed
		err = sendResponseInfo(conn, buildError(err.Error())) // Send error respone
		if err != nil {
			return // Exit upload
		}
	}
}
