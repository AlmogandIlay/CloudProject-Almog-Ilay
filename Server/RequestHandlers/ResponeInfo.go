package RequestHandlers

import (
	helper "CloudDrive/Helper"
	"encoding/json"
	"fmt"
	"net"
)

type ResponeType int

const (
	ErrorRespone ResponeType = 999
	ValidRespone ResponeType = 200
)

type ResponeInfo struct {
	Type       ResponeType
	Respone    string
	NewHandler *IRequestHandler
}

// Given responseinfo to client
type ClientResponeInfo struct {
	Type    ResponeType `json:"Type"`
	Respone string      `json:"Data"`
}

func buildRespone(respone string, handler *IRequestHandler) ResponeInfo {
	return ResponeInfo{Type: ValidRespone, Respone: respone, NewHandler: handler}
}

func buildError(response string, irequest IRequestHandler) ResponeInfo {
	return ResponeInfo{Type: ErrorRespone, Respone: response, NewHandler: &irequest}
}

func SendResponseInfo(conn *net.Conn, responseInfo ResponeInfo) error {
	clientResponeInfo := ClientResponeInfo{Type: responseInfo.Type, Respone: responseInfo.Respone}
	message, err := json.Marshal(clientResponeInfo)
	if err != nil {
		return err
	}
	fmt.Println("Response Json data is", string(message))
	return helper.SendData(conn, message)
}
