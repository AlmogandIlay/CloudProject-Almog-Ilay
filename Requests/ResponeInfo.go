package Requests

type ResponeType int

const (
	ErrorRespone ResponeType = 999
	ValidRespone ResponeType = 200
)

type ResponeInfo struct {
	Type    ResponeType `json:"Type"`
	Respone string      `json:"Data"`
}
