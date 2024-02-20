package Authentication

import (
	"client/ClientErrors"
	"client/Requests"
	"encoding/json"
	"fmt"
	"net"
)

const (
	signupArguments = 3
	loginArguments  = 2

	username_index = 0
	password_index = 1
	email_index    = 2
)

// function, argumentCount, arguments,

// Handles the sign up request
func HandleSignup(commandArguments []string, socket net.Conn) error {
	if len(commandArguments) != signupArguments {
		return &(ClientErrors.InvalidArgumentCountError{Arguments: uint8(len(commandArguments)), Expected: uint8(signupArguments)})
	}
	user := Signup(commandArguments[username_index], commandArguments[password_index], commandArguments[email_index])
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
		return &ClientErrors.JsonEncodeError{}
	}
	_, err = Requests.SendRequest(Requests.LoginRequest, request_data, socket)

	return err
}
