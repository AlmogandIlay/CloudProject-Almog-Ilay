package Requests

import (
	"Client/helper"
	"encoding/json"
	"fmt"
	"net"
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

func SendRequestInfo(request_info RequestInfo, socket net.Conn) (ResponeInfo, error) {
	requestBytes, err := json.Marshal(request_info)
	if err != nil {
		return ResponeInfo{}, fmt.Errorf("error when attempting to decode the data to be sent to the server.\nPlease send this info to the developers:\n%s", err.Error())
	}

	_, err = socket.Write(requestBytes)
	if err != nil {
		return ResponeInfo{}, fmt.Errorf("error when attempting to send the request to the server.\nPlease send this info to the developers:\n%s", err)
	}
	data, err := helper.ReciveData(socket, helper.DefaultBufferSize)
	response_info, err := helper.GetResponseInfo(data)

	return response_info, nil
}
