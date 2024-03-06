package RequestHandlers

import (
	helper "CloudDrive/Helper"
	"encoding/json"
	"net"
)

type ResponeType int

const (
	ErrorRespone  ResponeType = 999
	ValidRespone  ResponeType = 200
	OkayRespone   string      = "200: Okay" // Valid Response default data
	CDRespone     string      = "200: CurrentDirectory:"
	ChunksRespone string      = "200: ChunksSize:"
	SizeRespone   string      = " FileSize:"
)

type ResponeInfo struct {
	Type       ResponeType
	Respone    string
	NewHandler *IRequestHandler
}

// ResponseInfo for client
type ClientResponeInfo struct {
	Type    ResponeType `json:"Type"`
	Respone string      `json:"Data"`
}

// Create a valid response info
func buildRespone(respone string, handler *IRequestHandler) ResponeInfo {
	return ResponeInfo{Type: ValidRespone, Respone: respone, NewHandler: handler}
}

// Create an error response info
func buildError(response string, irequest IRequestHandler) ResponeInfo {
	return ResponeInfo{Type: ErrorRespone, Respone: response, NewHandler: &irequest}
}

// Send the ResponseInfo to the client
func SendResponseInfo(conn *net.Conn, responseInfo ResponeInfo) error {
	clientResponeInfo := ClientResponeInfo{Type: responseInfo.Type, Respone: responseInfo.Respone}
	message, err := json.Marshal(clientResponeInfo)
	if err != nil {
		return err
	}
	return helper.SendData(conn, message)
}
