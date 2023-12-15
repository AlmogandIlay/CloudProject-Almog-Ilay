package Respones

import "RequestHandlers"

type ResponeType int

const (
	LoginRespone  ResponeType = iota // 0
	signinRespone ResponeType = iota // 1
)

type ResponeInfo struct {
	messageCode ResponeType
	newHandler  *RequestHandlers.IRequestHandler
}
