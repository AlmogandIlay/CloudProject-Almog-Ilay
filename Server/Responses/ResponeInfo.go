package Responses

import (
	"CloudDrive/Server/RequestHandlers/Requests"
)

type ResponeType int

const (
	LoginRespone  ResponeType = iota // 0
	signinRespone ResponeType = iota // 1
)

type IRequestHandler interface {
	ValidRequest(info Requests.RequestInfo) bool
	HandleRequest(info Requests.RequestInfo) ResponeInfo
	Error() ResponeInfo
}

type ResponeInfo struct {
	messageCode ResponeType
	newHandler  IRequestHandler
}
