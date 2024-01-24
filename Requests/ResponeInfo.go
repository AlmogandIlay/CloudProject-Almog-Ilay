package Requests

type ResponeType int

const (
	ErrorRespone  ResponeType = iota // 0
	LoginRespone  ResponeType = iota // 1
	SignupRespone ResponeType = iota // 2
)

type ResponeInfo struct {
	Type    ResponeType `json:"Type"`
	Respone string      `json:"Data"`
}
