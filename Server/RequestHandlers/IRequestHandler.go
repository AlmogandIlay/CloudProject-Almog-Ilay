package requesthandlers

import "Requests"

type IRequestHandler interface {
	isValidRequest(info Requests.RequestInfo) bool
	handleRequest()
	error()
}
