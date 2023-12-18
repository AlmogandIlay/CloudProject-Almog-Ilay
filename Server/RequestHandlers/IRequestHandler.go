package RequestHandlers

import (
	"server/RequestHandlers/Requests"
)

type IRequestHandler interface {
	ValidRequest(info Requests.RequestInfo) bool
	HandleRequest(info Requests.RequestInfo) ResponeInfo
	Error(info Requests.RequestInfo) ResponeInfo
}
