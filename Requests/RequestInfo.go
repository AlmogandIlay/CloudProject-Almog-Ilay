package Requests

import (
	"client/ClientErrors"
	"client/Helper"
	"encoding/json"
	"fmt"
	"net"
)

type RequestType int

const (
	LoginRequest           RequestType = 101
	SignupRequest          RequestType = 102
	ChangeDirectoryRequest RequestType = 301
	CreateFileRequest      RequestType = 302
	CreateFolderRequest    RequestType = 303
	DeleteContentRequest   RequestType = 304
	RenameRequest          RequestType = 305
	ShowRequest            RequestType = 306
	MoveRequest            RequestType = 307
	GarbageRequest         RequestType = 308
	UploadFileRequest      RequestType = 401
	DownloadFileRequest    RequestType = 402
	UploadDirectoryRequest RequestType = 403
	DownloadDirRequest     RequestType = 404
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
		return ResponeInfo{}, &ClientErrors.JsonEncodeError{Err: err}
	}

	err = Helper.SendData(&socket, requestBytes)
	if err != nil {
		return ResponeInfo{}, err
	}

	data, err := Helper.ReciveData(&socket)
	if err != nil {
		return ResponeInfo{}, err
	}
	response_info, err := getResponseInfo(data)
	if err != nil {
		return ResponeInfo{}, err
	}

	return response_info, nil
}

// Handles the entire request-response cycle.
func SendRequest(requestType RequestType, request_data []byte, socket *net.Conn) (string, error) {
	request_info := BuildRequestInfo(requestType, request_data)
	response_info, err := SendRequestInfo(request_info, *socket) // sends a request and receives a response
	if err != nil {
		return "", err
	}
	if response_info.Type == ValidRespone { // If error caught in server side
		return response_info.Respone, nil
	} else {
		return "", fmt.Errorf(response_info.Respone)
	}
}
