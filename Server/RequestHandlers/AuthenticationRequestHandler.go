package RequestHandlers

import (
	helper "CloudDrive/Helper"
	"CloudDrive/Server/RequestHandlers/Requests"
)

type AuthenticationRequestHandler struct{}

func (loginHandler AuthenticationRequestHandler) HandleRequest(info Requests.RequestInfo) ResponeInfo {
	switch info.Type {
	case Requests.LoginRequest:
		return loginHandler.HandleLogin(info)
	case Requests.SignupRequest:
		return loginHandler.HandleSignup(info)
	default:
		return Error(info, IRequestHandler(&loginHandler))
	}

}

/*
Handle Login requests from client
*/
func (loginHandler *AuthenticationRequestHandler) HandleLogin(info Requests.RequestInfo) ResponeInfo {
	user := helper.GetEncodedUser(info.RequestData)
	login_manager := GetManager()

	err := login_manager.Login(user.Username, user.Password) // Attempt to perform a login request
	if err != nil {
		return buildError(err.Error(), loginHandler)
	}

	fileRequestHandler := FileRequestHandler{}
	var irequesthandler IRequestHandler = &fileRequestHandler
	return buildRespone("200: Okay", &irequesthandler) // Login request success (tdl: add handler)

}

func (loginHandler *AuthenticationRequestHandler) HandleSignup(info Requests.RequestInfo) ResponeInfo {
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

	fileRequestHandler := FileRequestHandler{}                    // Initialize file handler
	var irequestFileHandler IRequestHandler = &fileRequestHandler // convert the file handler to an interface
	return buildRespone("200: Okay", &irequestFileHandler)        // Signup request success

}
