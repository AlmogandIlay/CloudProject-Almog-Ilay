package RequestHandlers

import (
	"server/Requests"
)

type IRequestHandler interface {
	isValidRequest(info Requests.RequestInfo) bool
	handleRequest(info Requests.RequestInfo)
	error()
}
