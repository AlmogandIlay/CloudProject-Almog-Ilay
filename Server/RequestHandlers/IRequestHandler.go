package RequestHandlers

import (
	"CloudDrive/Server/RequestHandlers/Requests"
)

type IRequestHandler interface {
	ValidRequest(info Requests.RequestInfo) bool
	HandleRequest(info Requests.RequestInfo) ResponeInfo
}

func ChangeRequestHandler(response ResponeInfo) IRequestHandler {
	return *response.NewHandler

}

func Error(info Requests.RequestInfo, handler IRequestHandler) ResponeInfo {
	if info.Type == Requests.ErrorRequest { // If error request caught
		return buildError(string(info.RequestData), handler)
	}

	return buildError("Error: Not Exist.", handler) // Invalid request type

}
