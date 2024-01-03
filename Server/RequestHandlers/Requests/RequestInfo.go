package Requests

type RequestType int

const (
	LoginRequest  RequestType = iota // 0
	SignupRequest RequestType = iota // 1
)

type RequestInfo struct {
	Type        RequestType
	RequestData string
}
