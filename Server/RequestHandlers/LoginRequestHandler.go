package RequestHandlers

import (
	"server/RequestHandlers/Requests"
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

func (loginHandler *LoginRequestHandler) HandleLogin(info Requests.RequestInfo) ResponeInfo {

}

func (loginHandler *LoginRequestHandler) HandleSignin(info Requests.RequestInfo) ResponeInfo {
}
