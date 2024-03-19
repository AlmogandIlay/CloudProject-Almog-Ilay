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
	DeleteContentRequest   RequestType = 304
	RenameContentRequest   RequestType = 305
	ListContentsRequest    RequestType = 306
	MoveContentRequest     RequestType = 307
	GarbageRequest         RequestType = 308

	UploadFileRequest   RequestType = 401
	DownloadFileRequest RequestType = 402

	UploadDirectoryRequest RequestType = 403
	DownloadDirRequest     RequestType = 404

	StopTranmission RequestType = 501

	ErrorRequest RequestType = 999
)

type RequestInfo struct {
	Type        RequestType     `json:"Type"`
	RequestData json.RawMessage `json:"Data"`

	// Add Time data for log

}

// Recives request info json data from client socket and returns RequestInfo encoded struct with timeout flag
func ReciveRequestInfo(conn *net.Conn, timeoutFlag bool) (RequestInfo, error) {
	data, err := helper.ReciveData(conn, helper.DefaultBufferSize, timeoutFlag)

	if err != nil {
		return RequestInfo{}, err
	}

	var requestInfo RequestInfo

	err = json.Unmarshal(data, &requestInfo) // Encode json to struct
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
