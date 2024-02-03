package RequestHandlers

import (
	"CloudDrive/FileSystem"
	"CloudDrive/Server/RequestHandlers/Requests"
	"fmt"
)

type FileRequestHandler struct{}

func (filehandler FileRequestHandler) HandleRequest(info Requests.RequestInfo) ResponeInfo {
	switch info.Type {
	case Requests.ChangeDirectoryRequest:
		return filehandler.handleChangeDirectory(info)
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

	}
	return ResponeInfo{}
}

func (filehandler *FileRequestHandler) handleChangeDirectory(info Requests.RequestInfo) ResponeInfo {
	path := string(info.RequestData)
	err := FileSystem.ChangeDirectory(path)
	fmt.Println(string(info.RequestData))
	return ResponeInfo{}
}
