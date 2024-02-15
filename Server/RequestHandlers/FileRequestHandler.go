package RequestHandlers

import (
	"CloudDrive/FileSystem"
	helper "CloudDrive/Helper"
	"CloudDrive/Server/RequestHandlers/Requests"
	"strings"
)

const (
	fileName    = 0
	newFileName = 1
)

type FileRequestHandler struct{}

// Handle the file system type requests
func (filehandler FileRequestHandler) HandleRequest(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
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
	case Requests.RenameFileRequest:
		return filehandler.handleRenameFile(info, loggedUser)
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

// Handle Rename File requests from client
func (filehandler *FileRequestHandler) handleRenameFile(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	command := Requests.ParseDataToString(info.RequestData)
	arguments := strings.Fields(command)
	err := loggedUser.RenameFile(arguments[fileName], arguments[newFileName])
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler())
}

// Handle List Contents (ls) requests from client
func (filehandler *FileRequestHandler) handleListContents(loggeduser *FileSystem.LoggedUser) ResponeInfo {
	list, err := loggeduser.ListContents()
	if err != nil {
		return buildError(err.Error(), IRequestHandler(filehandler))
	}

	return buildRespone(list, CreateFileRequestHandler())
}
