package Requests

type RequestType int

const (
	LoginRequest  RequestType = 101
	SignupRequest RequestType = 102
	// ...
	ErrorRequest RequestType = 999
)

type RequestInfo struct {
	Type        RequestType `json:"Type"`
	RequestData []byte      `json:"Data"`

	// Add Time data for log

}
