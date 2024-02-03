package RequestHandlers

import (
	"CloudDrive/FileSystem"
	"CloudDrive/Server/RequestHandlers/Requests"
)

type FileRequestHandler struct{}

func (filehandler FileRequestHandler) HandleRequest(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	switch info.Type {
	case Requests.ChangeDirectoryRequest:
		return filehandler.handleChangeDirectory(info, loggedUser)
	case Requests.CreateFileRequest:
		// 		return loginHandler.HandleSignup(info)
		// 	case Requests.CreateFolderRequest:
		// 		// TODO
		// 	case Requests.DeleteFileRequest:
		// 		// TODO
		// 	case Requests.DeleteFolderRequest:
		// 		// TODO
		// 	default:
		// 		return Error(info, IRequestHandler(filehandler))
		// 	}

		// }

		// func (fileHanlder *FileRequestHandler) HandleChangeDirectory(info Requests.RequestInfo) ResponeInfo {

	}
	return ResponeInfo{}
}

// Handle cd (Change Directory) requests from client
func (filehandler *FileRequestHandler) handleChangeDirectory(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	path := string(info.RequestData)
	err := loggedUser.ChangeDirectory(path)
	if err != nil {
		buildError(err.Error(), IRequestHandler(filehandler))
	}

	fileRequestHandler := FileRequestHandler{}                    // Initialize file handler
	var irequestFileHandler IRequestHandler = &fileRequestHandler // convert the file handler to an interface
	return buildRespone(OkayRespone, &irequestFileHandler)
}
