package RequestHandlers

import (
	"CloudDrive/FileSystem"
	helper "CloudDrive/Helper"
	"CloudDrive/Server/RequestHandlers/Requests"
	"net"
	"strconv"
	"strings"
)

const (
	pathContentName    = 1
	newPathContentName = 3
)

type FileRequestHandler struct{}

// Handle the file system type requests
func (filehandler FileRequestHandler) HandleRequest(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser, fileTransferListener *net.Listener) ResponeInfo {
	switch info.Type {
	case Requests.ChangeDirectoryRequest:

		return filehandler.handleChangeDirectory(info, loggedUser)
	case Requests.CreateFileRequest:
		return filehandler.handleCreateFile(info, loggedUser)
	case Requests.CreateFolderRequest:
		return filehandler.handleCreateFolder(info, loggedUser)
	case Requests.DeleteContentRequest:
		return filehandler.handleDeleteContent(info, loggedUser)
	case Requests.RenameContentRequest:
		return filehandler.handleRenameContent(info, loggedUser)
	case Requests.ListContentsRequest:
		return filehandler.handleListContents(info, loggedUser)
	case Requests.MoveContentRequest:
		return filehandler.handleMoveContent(info, loggedUser)
	case Requests.UploadFileRequest:
		return filehandler.handleUploadFile(info, loggedUser, fileTransferListener)
	case Requests.DownloadFileRequest:
		return filehandler.handleDownloadFile(info, loggedUser, fileTransferListener)
	case Requests.GarbageRequest:
		return filehandler.handleGarbage(loggedUser)
	case Requests.UploadDirectoryRequest:
		return filehandler.handleUploadDirectory(info, loggedUser, fileTransferListener)
	case Requests.DownloadDirRequest:
		return filehandler.handleDownloadDir(info, loggedUser, fileTransferListener)
	default:
		return Error(info, IRequestHandler(&filehandler))
	}
}

// Handle the garbage command, chage the current directory to garbage directory
func (filehandler *FileRequestHandler) handleGarbage(loggedUser *FileSystem.LoggedUser) ResponeInfo {
	path, err := loggedUser.GarbageChangeDirectory() // Changes current directory to Garbage path
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(CDRespone+path, CreateFileRequestHandler())
}

// Handle cd (Change Directory) requests from client
func (filehandler *FileRequestHandler) handleChangeDirectory(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	rawData := Requests.ParseDataToString(info.RequestData) // Convert the RequestInfo.Data to raw string
	requestPath := helper.ConvertRawJsonToData(rawData)     // Fixes the raw string to path string
	path, err := loggedUser.ChangeDirectory(requestPath)    // Changes the current directory to the given path string
	if err != nil {                                         // If error has occured
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(CDRespone+path, CreateFileRequestHandler())
}

// Handle Create File requests from client
func (filehandler *FileRequestHandler) handleCreateFile(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	rawData := Requests.ParseDataToString(info.RequestData) // Convert the RequestInfo.Data to raw string
	file := helper.ConvertRawJsonToData(rawData)            // Fixes the raw string to path string
	err := loggedUser.CreateFile(file)                      // Creates a file with the given path string
	if err != nil {                                         // If error has occured
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler())
}

// Handle Create Folder requests from client
func (filehandler *FileRequestHandler) handleCreateFolder(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	rawData := Requests.ParseDataToString(info.RequestData) // Convert the RequestInfo.Data to raw string
	folderName := helper.ConvertRawJsonToData(rawData)      // Fixes the raw string to path string
	err := loggedUser.CreateFolder(folderName)              // Creates a folder with the given path string
	if err != nil {                                         // If error has occured
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler())
}

// Handle Delete Content requests from client
func (filehandler *FileRequestHandler) handleDeleteContent(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	rawData := Requests.ParseDataToString(info.RequestData) // Convert the RequestInfo.Data to raw string
	content := helper.ConvertRawJsonToData(rawData)         // Fixes the raw string to path string
	err := loggedUser.RemoveContent(content)                // Removes a content with the given path string
	if err != nil {                                         // If error has occured
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler())
}

// Handle Rename content (file and folders) requests from client
func (filehandler *FileRequestHandler) handleRenameContent(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	command := Requests.ParseDataToString(info.RequestData)                                    // Convert the RequestInfo.Data to raw string
	data := helper.ConvertRawJsonToData(command)                                               // Fixes the raw string to path string
	arguments := strings.Split(data, "'")                                                      // Saves the rename arguments in a slice
	err := loggedUser.RenameContent(arguments[pathContentName], arguments[newPathContentName]) // Renames a given content name with the given new content name
	if err != nil {                                                                            // If error has occured
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler())
}

// Handle List Contents (ls) requests from client
func (filehandler *FileRequestHandler) handleListContents(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	command := Requests.ParseDataToString(info.RequestData) // Convert request data to command string
	var path string
	if command != "null" { // If path has been specified
		path = helper.ConvertRawJsonToData(command) // Save specified path
	}
	list, err := loggedUser.ListContents(path) // List the contents inside the given path
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(list, CreateFileRequestHandler())
}

// Handle Move content (files and folders) requests from client
func (filehandler *FileRequestHandler) handleMoveContent(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	command := Requests.ParseDataToString(info.RequestData)                                  // Convert the RequestInfo.Data to raw string
	data := helper.ConvertRawJsonToData(command)                                             // Fixes the raw string to path string
	arguments := strings.Split(data, "'")                                                    // Saves the move arguments in a slice
	err := loggedUser.MoveContent(arguments[pathContentName], arguments[newPathContentName]) // Moves a given content name to the given new directory
	if err != nil {                                                                          // If error has occured
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler())
}

// Handle Upload file requests from client
func (filehandler *FileRequestHandler) handleUploadFile(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser, uploadListener *net.Listener) ResponeInfo {
	file, err := FileSystem.ParseDataToContent(info.RequestData) // Parse RequestInfo.Data to Content struct
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	chunksSize, err := loggedUser.UploadFile(file, uploadListener) // Sends to upload file request
	if err != nil {                                                // If the file to upload is not valid
		return buildError(err.Error(), IRequestHandler(filehandler))
	}
	// The file to upload is valid

	return buildRespone(ChunksRespone+strconv.FormatUint(uint64(chunksSize), 10), CreateFileRequestHandler()) // Sends chunks size of the given file
}

// Handle Download file requests from client
func (filehandler *FileRequestHandler) handleDownloadFile(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser, downloadListener *net.Listener) ResponeInfo {
	rawData := Requests.ParseDataToString(info.RequestData)                          // Convert the RequestInfo.Data to raw string
	filename := helper.ConvertRawJsonToData(rawData)                                 // Fixes the raw string to path string
	chunksSize, fileSize, err := loggedUser.DownloadFile(filename, downloadListener) // Sends to download file request
	if err != nil {                                                                  // If the file to download is not valid
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(ChunksRespone+strconv.FormatUint(uint64(chunksSize), 10)+SizeRespone+strconv.FormatUint(uint64(fileSize), 10), CreateFileRequestHandler()) // Sends chunk size and the file size of the given file.
}

// Handle Upload directory requests from client
func (filehandler *FileRequestHandler) handleUploadDirectory(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser, uploadListener *net.Listener) ResponeInfo {
	dir, err := FileSystem.ParseDataToContent(info.RequestData) // Parse RequestInfo.Data to Content struct
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}
	err = loggedUser.UploadDirectory(dir, uploadListener) // Sends to upload directory request
	if err != nil {                                       // If the directory to upload is not valid
		return buildError(err.Error(), IRequestHandler(filehandler))
	}
	return buildRespone(OkayRespone, CreateFileRequestHandler()) // Send Directory is valid respone

}

func (filehandler *FileRequestHandler) handleDownloadDir(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser, downloadListener *net.Listener) ResponeInfo {
	rawData := Requests.ParseDataToString(info.RequestData)

	path := helper.ConvertRawJsonToData(rawData) // Fixes the raw string to path string

	err := loggedUser.DownloadDirectory(path, downloadListener)
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler())

}
