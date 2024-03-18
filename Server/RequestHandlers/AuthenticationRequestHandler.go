package RequestHandlers

import (
	"CloudDrive/FileSystem"
	helper "CloudDrive/Helper"
	"CloudDrive/Server/RequestHandlers/Requests"
	"net"
)

type AuthenticationRequestHandler struct{}

// Handle Authentication type requests
func (loginHandler AuthenticationRequestHandler) HandleRequest(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser, uploadListener *net.Listener) ResponeInfo {
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
	user := helper.GetEncodedUser(info.RequestData) // Encode RequestInfo to User struct
	login_manager := GetManager()                   // Recives a pointer to LoginManager API

	err := login_manager.Login(user.Username, user.Password) // Attempt to perform a login request
	if err != nil {
		return buildError(err.Error(), loginHandler)
	}
	// User has successfully signed in

	*loggedUser, err = GetLoggedUser(info) // Recives LoggedUser API Instance and saves it to the loggedUser parameter
	if err != nil {                        // If error has occured
		return buildError(err.Error(), loginHandler)
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler()) // Login request has success

}

// Handle Signup requests from client
func (loginHandler *AuthenticationRequestHandler) handleSignup(info Requests.RequestInfo, loggedUser *FileSystem.LoggedUser) ResponeInfo {
	user := helper.GetEncodedUser(info.RequestData) // Encode RequestInfo to User struct
	login_manager := GetManager()                   // Recives a pointer to LoginManager API

	errs := login_manager.Signup(user.Username, user.Password, user.Email) // Attempt to perform a signup request
	if len(errs) > 0 {                                                     // If errors occured in the signup process
		var errors string = ""
		for _, err := range errs { // Save all errors data in one string
			errors += "* " + err.Error() + "\n"
		}
		return buildError(errors, loginHandler)
	}
	// User has successfully registered

	var err error
	*loggedUser, err = GetLoggedUser(info) // Initializes LoggedUser API Instance and saves it to the loggedUser parameter
	if err != nil {                        // If error has occured
		return buildError(err.Error(), loginHandler)
	}

	return buildRespone(OkayRespone, CreateFileRequestHandler()) // Signup request has success

}
