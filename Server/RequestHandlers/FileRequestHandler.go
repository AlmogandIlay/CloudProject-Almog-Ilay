package RequestHandlers

import (
	"CloudDrive/Server/RequestHandlers/Requests"
)

type FileRequestHandler struct{}

func (loginHandler *FileRequestHandler) HandleRequest(info Requests.RequestInfo) ResponeInfo {
	switch info.Type {
	case Requests.ChangeDirectoryRequest:
		return
	case Requests.CreateFileRequest:
		return loginHandler.HandleSignup(info)
	case Requests.CreateFolderRequest:
		// TODO
	case Requests.DeleteFileRequest:
		// TODO
	case Requests.DeleteFolderRequest:
		// TODO
	default:
		return loginHandler.HandleError(info)
	}

}

func (fileHanlder *FileRequestHandler) HandleChangeDirectory(info Requests.RequestInfo) ResponeInfo {

}
