package Requests

type RequestType int

const (
	LoginRequest  RequestType = iota // 0
	signinRequest RequestType = iota // 1
)

type RequestInfo struct {
	messageCode RequestType
	request     string
}
