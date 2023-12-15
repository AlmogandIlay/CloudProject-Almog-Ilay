package RequestHandlers

import (
	"Requests"
)

type IRequestHandler interface {
	isValidRequest(info Requests.RequestInfo) bool
	handleRequest(info Requests.RequestInfo)
	error()
}
