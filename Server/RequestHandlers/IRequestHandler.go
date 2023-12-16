package RequestHandlers

import (
	"server/Requests"
	"server/Responses"
)

type IRequestHandler interface {
	ValidRequest(info Requests.RequestInfo) bool
	HandleRequest(info Requests.RequestInfo) Responses.ResponeInfo
	Error() Responses.ResponeInfo
}
