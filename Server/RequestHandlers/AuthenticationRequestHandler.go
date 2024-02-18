package RequestHandlers

import (
	"CloudDrive/FileSystem"
	helper "CloudDrive/Helper"
	"CloudDrive/Server/RequestHandlers/Requests"
	"fmt"
	"net"
)

type AuthenticationRequestHandler struct{}

// Handle Authentication type requests
func (loginHandler AuthenticationRequestHandler) HandleRequest(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser, conn *net.Conn) ResponeInfo {
	switch info.Type {
	case Requests.LoginRequest:
		return loginHandler.handleLogin(info, loggedUser)
	case Requests.SignupRequest:
		return loginHandler.handleSignup(info, loggedUser)
	default:
		return Error(info, IRequestHandler(&loginHandler))
	}

}

/*
Handle Login requests from client
*/
func (loginHandler *AuthenticationRequestHandler) handleLogin(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	user := helper.GetEncodedUser(info.RequestData)
	login_manager := GetManager()

	err := login_manager.Login(user.Username, user.Password) // Attempt to perform a login request
	if err != nil {
		return buildError(err.Error(), loginHandler)
	}

	*loggedUser, err = GetLoggedUser(info)
	if err != nil {
		return buildError(err.Error(), loginHandler)
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler()) // Login request success (tdl: add handler)

}

// Handle Signup requests from client
func (loginHandler *AuthenticationRequestHandler) handleSignup(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	user := helper.GetEncodedUser(info.RequestData)
	login_manager := GetManager()

	errs := login_manager.Signup(user.Username, user.Password, user.Email) // Attempt to perform a signup request
	if len(errs) > 0 {
		var errors string = ""
		for _, err := range errs { // Save all errors in string
			errors += "* " + err.Error() + "\n"
		}
		return buildError(errors, loginHandler)
	}

	var err error
	*loggedUser, err = GetLoggedUser(info)
	if err != nil {
		fmt.Println(err.Error())
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler()) // Signup request success

}
