package RequestHandlers

import (
	"CloudDrive/FileSystem"
	"CloudDrive/Server/RequestHandlers/Requests"
)

type IRequestHandler interface {
	HandleRequest(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo
}

func UpdateRequestHandler(response ResponeInfo) IRequestHandler {
	return *response.NewHandler

}

func Error(info Requests.RequestInfo, handler IRequestHandler) ResponeInfo {
	if info.Type == Requests.ErrorRequest { // If error request caught
		return buildError(string(info.RequestData), handler)
	}

	return buildError("Error: Unknown Command.", handler) // Invalid request type

}

func CreateFileRequestHandler() *IRequestHandler {
	FileRequestHandler := FileRequestHandler{}                    // Initialize file handler
	var irequestFileHandler IRequestHandler = &FileRequestHandler // convert the file handler to an interface
	return &irequestFileHandler
}
