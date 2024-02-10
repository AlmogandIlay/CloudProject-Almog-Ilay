package Authentication

import (
	"client/Requests"
	"encoding/json"
	"fmt"
	"net"
)

const (
	username_index = 0
	password_index = 1
	email_index    = 2
)

// function, argumentCount, arguments,

// Handles the sign up request
func HandleSignup(command_arguments []string, socket net.Conn) error {
	if len(command_arguments) != 3 {
		return fmt.Errorf("incorrect number of arguments.\nPlease try again")
	}
	user := Signup(command_arguments[username_index], command_arguments[password_index], command_arguments[email_index])
	request_data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Error when attempting to encode the data to be sent to the server.\nPlease send this info to the developers: %s", err.Error()))
	}
	_, err = Requests.SendRequest(Requests.SignupRequest, request_data, socket)

	return err

}

func HandleSignIn(command_arguments []string, socket net.Conn) error {
	if len(command_arguments) != 2 {
		return fmt.Errorf("incorrect number of arguments.\nPlease try again")
	}
	user := Signin(command_arguments[username_index], command_arguments[password_index])
	request_data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Error when attempting to encode the data to be sent to the server.\nPlease send this info to the developers: %s", err.Error()))
	}
	_, err = Requests.SendRequest(Requests.LoginRequest, request_data, socket)

	return err
}
