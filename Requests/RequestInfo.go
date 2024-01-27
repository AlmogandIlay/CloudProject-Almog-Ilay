package Requests

import (
	"client/Helper"
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

	err = Helper.SendData(&socket, requestBytes)
	if err != nil {
		return ResponeInfo{}, err
	}

	data, err := Helper.ReciveData(&socket, Helper.DefaultBufferSize)
	if err != nil {
		return ResponeInfo{}, err
	}

	response_info, err := GetResponseInfo(data)
	if err != nil {
		return ResponeInfo{}, err
	}

	return response_info, nil
}

// Handles the entire request-response cycle.
func SendRequest(request_data []byte, socket net.Conn) error {
	request_info := BuildRequestInfo(SignupRequest, request_data)
	response_info, err := SendRequestInfo(request_info, socket) // sends a request and receives a response
	if err != nil {
		return err
	}
	if response_info.Type == ErrorRespone { // If error caught in server side
		return fmt.Errorf(response_info.Respone)
	} else if response_info.Type == ValidRespone {
		return nil
	}
	return fmt.Errorf(response_info.Respone)
}
