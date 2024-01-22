package Requests

import (
	"encoding/json"
	"fmt"
)

type RequestType int

const (
	LoginRequest  RequestType = 101
	SignupRequest RequestType = 102
	// ...
)

type RequestInfo struct {
	Type        RequestType     `json:"Type"`
	RequestData json.RawMessage `json:"Data"`

	// Add Time data for log

}

// BuildRequestInfo creates a new RequestInfo struct with the given request type and data.

func BuildRequestInfo(request_type RequestType, request_data json.RawMessage) RequestInfo {
	return RequestInfo{
		Type:        request_type,
		RequestData: request_data,
	}
}

func SendRequestInfo(request_info RequestInfo) ResponeInfo {
	requestBytes, err := json.Marshal(request_info)
	if err != nil {
		panic(fmt.Sprintf("Error when attempting to decode the data to be sent to the server.\nPlease send this info to the developers:\n%s", err.Error()))
	}
	fmt.Println(requestBytes)
	return ResponeInfo{}

	//  _, err = cli.socket.Write(requestBytes)
}
