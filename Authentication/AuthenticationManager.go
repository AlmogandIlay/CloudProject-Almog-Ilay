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

	request_info := Requests.BuildRequestInfo(Requests.SignupRequest, request_data)
	response_info, err := Requests.SendRequestInfo(request_info, socket)
	if err != nil {
		return err
	}
	if response_info.Type == Requests.ErrorRespone { // If error caught in server side
		return fmt.Errorf(response_info.Respone)
	}

	if response_info.Type == Requests.ValidRespone {
		fmt.Println("Signup successful!")
	}

	return nil

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
	request_info := Requests.BuildRequestInfo(Requests.LoginRequest, request_data)
	response_info, err := Requests.SendRequestInfo(request_info, socket)
	if err != nil {
		return err
	}
	fmt.Println(response_info.Type, response_info.Respone)
	return nil
}
