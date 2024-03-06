package Requests

import (
	helper "CloudDrive/Helper"
	"encoding/json"
	"net"
	"strings"
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
	RenameContentRequest   RequestType = 306
	ListContentsRequest    RequestType = 307
	MoveContentRequest     RequestType = 308
	GarbageRequest         RequestType = 309
	UploadFileRequest      RequestType = 401
	DownloadFileRequest    RequestType = 402

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
	if err != nil {
		return RequestInfo{
			Type:        ErrorRequest,
			RequestData: []byte("Error: One or more of the fields are being used wrong."),
		}, nil

	}

	return requestInfo, nil

}

// Convert raw json data to string
func ParseDataToString(data json.RawMessage) string {
	fixData := strings.Replace(string(data), `\\`, `\`, -1)
	return fixData
}
