package RequestHandlers

import (
	"CloudDrive/Server/RequestHandlers/Requests"
	"CloudDrive/authentication"
	"encoding/json"
	"fmt"
)

type LoginRequestHandler struct{}

func (loginHandler *LoginRequestHandler) ValidRequest(info Requests.RequestInfo) bool {
	return info.Type == Requests.LoginRequest || info.Type == Requests.SigninRequest
}

func (loginHandler *LoginRequestHandler) HandleRequest(info Requests.RequestInfo) ResponeInfo {
	switch info.Type {
	case Requests.LoginRequest:
		return loginHandler.HandleLogin(info)
	case Requests.SigninRequest:
		return loginHandler.HandleSignin(info)
	default:
		return loginHandler.Error(info)
	}

}

func (loginHandler *LoginRequestHandler) Error(info Requests.RequestInfo) ResponeInfo {
	var respone ResponeInfo

	respone.NewHandler = nil

	//respone.Respone = some_func(info.Request) // a function from message to protocol

	return respone
}

/*
Handle Login requests from client
*/
func (loginHandler *LoginRequestHandler) HandleLogin(info Requests.RequestInfo) ResponeInfo { // add error to handles?
	var user authentication.User
	err := json.Unmarshal([]byte(info.RequestData), &user) // Json decoding
	if err != nil {
		fmt.Println(err.Error())
		return ResponeInfo{}
	}

	manager, err := GetManager() // Accessing Login Manager
	if err != nil {
		return ResponeInfo{}
	}

	err = manager.Login(user.Username, user.Password) // Attempting to perform a login request
	if err != nil {
		return ResponeInfo{}
	}

	return buildRespone("200: Okay", nil) // Login request success

}

func (loginHandler *LoginRequestHandler) HandleSignin(info Requests.RequestInfo) ResponeInfo {
	return ResponeInfo{}
}
