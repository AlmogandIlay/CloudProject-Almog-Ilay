package Requests

type RequestType int

const (
	LoginRequest  RequestType = iota // 0
	SignupRequest RequestType = iota // 1
	// ...
	ErrorRequest RequestType = 999
)

type RequestInfo struct {
	Type        RequestType `json:"Type"`
	RequestData []byte      `json:"Data"`

	// Add Time data for log

}
