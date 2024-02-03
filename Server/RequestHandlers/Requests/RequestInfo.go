package Requests

import (
	helper "CloudDrive/Helper"
	"encoding/json"
	"fmt"
	"net"
)

type RequestType int

const (
	// Authorization requests
	LoginRequest  RequestType = 101
	SignupRequest RequestType = 102

	// File requests
	ChangeDirectoryRequest RequestType = 301
	CreateFileRequest      RequestType = 302
	CreateFolderRequest    RequestType = 303
	DeleteFileRequest      RequestType = 304
	DeleteFolderRequest    RequestType = 305

	ErrorRequest RequestType = 999
)

type RequestInfo struct {
	Type        RequestType     `json:"Type"`
	RequestData json.RawMessage `json:"Data"`

	// Add Time data for log

}

// Recives data from client socket and returns RequestInfo
func ReciveRequestInfo(conn *net.Conn) (RequestInfo, error) {
	data, err := helper.ReciveData(conn, helper.DefaultBufferSize)

	if err != nil {
		return RequestInfo{}, err
	}

	var requestInfo RequestInfo

	err = json.Unmarshal(data, &requestInfo)
	fmt.Println("Request Json data is ", string(requestInfo.RequestData))
	if err != nil {
		return RequestInfo{
			Type:        ErrorRequest,
			RequestData: []byte("Error: One or more of the fields are being used wrong."),
		}, nil

	}

	return requestInfo, nil

}