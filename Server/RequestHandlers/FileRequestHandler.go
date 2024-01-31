package RequestHandlers

import (
	"CloudDrive/Server/RequestHandlers/Requests"
)

type FileRequestHandler struct{}

func (filehandler *FileRequestHandler) ValidRequest(info Requests.RequestInfo) bool {
	return info.Type == Requests.ChangeDirectoryRequest || info.Type == Requests.CreateFolderRequest || info.Type == Requests.CreateFileRequest || info.Type == Requests.DeleteFileRequest || info.Type == Requests.DeleteFolderRequest
}

func (filehandler *FileRequestHandler) HandleRequest(info Requests.RequestInfo) ResponeInfo {
	// 	switch info.Type {
	// 	case Requests.ChangeDirectoryRequest:
	// 		return
	// 	case Requests.CreateFileRequest:
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
	return buildRespone("", nil)
}
