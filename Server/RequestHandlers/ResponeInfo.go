package RequestHandlers

import (
	helper "CloudDrive/Helper"
	"encoding/json"
	"fmt"
	"net"
)

type ResponeType int

const (
	ErrorRespone  ResponeType = iota // 0
	LoginRespone  ResponeType = iota // 1
	SignupRespone ResponeType = iota // 2
)

type ResponeInfo struct {
	Respone    string
	NewHandler *IRequestHandler
}

func buildRespone(respone string, handler *IRequestHandler) ResponeInfo {
	return ResponeInfo{Respone: respone, NewHandler: handler}
}

func buildError(response string) ResponeInfo {
	return ResponeInfo{Respone: response, NewHandler: nil}
}

func SendResponseInfo(conn *net.Conn, responseInfo ResponeInfo) error {
	message, err := json.Marshal(responseInfo)
	if err != nil {
		return err
	}
	fmt.Println("Response Json data is ", string(message))
	return helper.SendData(conn, message)
}
