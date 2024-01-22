package Authentication

import (
	"client/Requests"
	"encoding/json"
	"fmt"
)

const (
	username_index = 0
	password_index = 1
	email_index    = 2
)

// Handles the sign up request
func HandleSignup(command_arguments []string) error {
	user := Signup(command_arguments[username_index], command_arguments[password_index], command_arguments[email_index])
	request_data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Error when attempting to encode the data to be sent to the server.\nPlease send this info to the developers: %s", err.Error()))
	}

	request_info := Requests.BuildRequestInfo(Requests.SignupRequest, request_data)
	fmt.Println(request_info)
	return nil
	//Requests.SendRequest(request_info)
}
