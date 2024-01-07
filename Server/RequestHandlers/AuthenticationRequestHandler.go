package RequestHandlers

import (
	"CloudDrive/Server/RequestHandlers/Requests"
	"CloudDrive/authentication"
	"encoding/json"
)

type AuthenticationRequestHandler struct{}

func (loginHandler *AuthenticationRequestHandler) ValidRequest(info Requests.RequestInfo) bool {
	return info.Type == Requests.LoginRequest || info.Type == Requests.SignupRequest
}

func (loginHandler *AuthenticationRequestHandler) HandleRequest(info Requests.RequestInfo) ResponeInfo {
	switch info.Type {
	case Requests.LoginRequest:
		return loginHandler.HandleLogin(info)
	case Requests.SignupRequest:
		return loginHandler.HandleSignup(info)
	default:
		return loginHandler.HandleError(info)
	}

}

func (loginHandler *AuthenticationRequestHandler) HandleError(info Requests.RequestInfo) ResponeInfo {
	if info.Type == Requests.ErrorRequest { // If error request caught
		return buildError(string(info.RequestData))
	}

	return buildError("Error: Not Exist.") // Invalid request type

}

/*
Handle Login requests from client
*/
func (loginHandler *AuthenticationRequestHandler) HandleLogin(info Requests.RequestInfo) ResponeInfo {
	var user authentication.User

	json.Unmarshal([]byte(info.RequestData), &user) // Json decoding
	login_manager := GetManager()

	err := login_manager.Login(user.Username, user.Password) // Attempt to perform a login request
	if err != nil {
		return buildError(err.Error())
	}

	return buildRespone("200: Okay", nil) // Login request success (tdl: add handler)

}

func (loginHandler *AuthenticationRequestHandler) HandleSignup(info Requests.RequestInfo) ResponeInfo {
	var user authentication.User

	json.Unmarshal([]byte(info.RequestData), &user) // Json decoding
	login_manager := GetManager()

	errs := login_manager.Signup(user.Username, user.Password, user.Email) // Attempt to perform a signup request
	if len(errs) > 0 {
		var errors string = ""
		for _, err := range errs { // Save all errors in string
			errors += "* " + err.Error() + "\n"
		}
		return buildError(errors)
	}

	return buildRespone("200: Okay", nil) // Signup request success (tdl: add handler)

}
