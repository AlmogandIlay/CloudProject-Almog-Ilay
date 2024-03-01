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
	case Requests.DeleteFileRequest:
		return filehandler.handleDeleteFile(info, loggedUser)
	case Requests.DeleteFolderRequest:
		return filehandler.handleDeleteFolder(info, loggedUser)
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
	default:
		return Error(info, IRequestHandler(&filehandler))
	}
}

// Handle cd (Change Directory) requests from client
func (filehandler *FileRequestHandler) handleChangeDirectory(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	rawData := Requests.ParseDataToString(info.RequestData)
	requestPath := helper.ConvertRawJsonToData(rawData)
	path, err := loggedUser.ChangeDirectory(requestPath)
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(CDRespone+path, CreateFileRequestHandler())
}

// Handle Create File requests from client
func (filehandler *FileRequestHandler) handleCreateFile(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	rawData := Requests.ParseDataToString(info.RequestData)
	file := helper.ConvertRawJsonToData(rawData)
	err := loggedUser.CreateFile(file)
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler())
}

// Handle Create Folder requests from client
func (filehandler *FileRequestHandler) handleCreateFolder(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	rawData := Requests.ParseDataToString(info.RequestData)
	folderName := helper.ConvertRawJsonToData(rawData)
	err := loggedUser.CreateFolder(folderName)
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler())
}

// Handle Delete File requests from client
func (filehandler *FileRequestHandler) handleDeleteFile(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	rawData := Requests.ParseDataToString(info.RequestData)
	file := helper.ConvertRawJsonToData(rawData)
	err := loggedUser.RemoveFile(file)
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler())
}

// Handle Delete Folder requests from client
func (filehandler *FileRequestHandler) handleDeleteFolder(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	rawData := Requests.ParseDataToString(info.RequestData)
	folderName := helper.ConvertRawJsonToData(rawData)
	err := loggedUser.RemoveFolder(folderName)
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler())
}

// Handle Rename content (file and folders) requests from client
func (filehandler *FileRequestHandler) handleRenameContent(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	command := Requests.ParseDataToString(info.RequestData)
	data := helper.ConvertRawJsonToData(command)
	arguments := strings.Split(data, "'")
	err := loggedUser.RenameContent(arguments[pathContentName], arguments[newPathContentName])
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler())
}

// Handle List Contents (ls) requests from client
func (filehandler *FileRequestHandler) handleListContents(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	path := ""
	if string(info.RequestData) != "null" { // If given path has been specified
		path = helper.ConvertRawJsonToData(string(info.RequestData))
	}
	list, err := loggedUser.ListContents(path)
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(list, CreateFileRequestHandler())
}

// Handle Move content (files and folders) requests from client
func (filehandler *FileRequestHandler) handleMoveContent(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	command := Requests.ParseDataToString(info.RequestData)
	data := helper.ConvertRawJsonToData(command)
	arguments := strings.Split(data, "'")
	err := loggedUser.MoveContent(arguments[pathContentName], arguments[newPathContentName])
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler())
}

// Handle Upload file requests from client
func (filehandler *FileRequestHandler) handleUploadFile(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser, uploadListener *net.Listener) ResponeInfo {
	file, err := FileSystem.ParseDataToFile(info.RequestData)
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	chunksSize, err := loggedUser.UploadFile(file, uploadListener)
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(ChunksRespone+strconv.FormatUint(uint64(chunksSize), 10), CreateFileRequestHandler()) // Send chunks size
}

func (filehandler *FileRequestHandler) handleDownloadFile(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser, downloadListener *net.Listener) ResponeInfo {
	rawData := Requests.ParseDataToString(info.RequestData)
	filename := helper.ConvertRawJsonToData(rawData)
	chunksSize, err := loggedUser.DownloadFile(filename, downloadListener)
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(ChunksRespone+strconv.FormatUint(uint64(chunksSize), 10), CreateFileRequestHandler())
}
