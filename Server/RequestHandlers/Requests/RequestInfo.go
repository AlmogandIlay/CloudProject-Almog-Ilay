package Requests

type RequestType int

const (
	LoginRequest  RequestType = iota // 0
	SigninRequest RequestType = iota // 1
)

type RequestInfo struct {
	Type        RequestType
	RequestData string
}
