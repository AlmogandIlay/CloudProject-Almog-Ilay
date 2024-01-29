package Requests

import "encoding/json"

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
