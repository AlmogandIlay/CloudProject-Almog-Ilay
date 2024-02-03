package RequestHandlers

import (
	"CloudDrive/FileSystem"
	helper "CloudDrive/Helper"
	"CloudDrive/Server/RequestHandlers/Requests"
	"CloudDrive/authentication"
)

var (
	Manager *authentication.AuthenticationManager
)

func InitializeIdentifyManagerFactory() (*authentication.AuthenticationManager, error) {
	var err error
	Manager, err = authentication.InitializeAuthenticationManager()
	if err != nil {
		Manager = &authentication.AuthenticationManager{}
		return Manager, err
	}

	return Manager, nil
}

func GetManager() *authentication.AuthenticationManager {
	return Manager
}

func GetLoggedUser(requestInfo Requests.RequestInfo, responeInfo ResponeInfo) (*FileSystem.LoggedUser, bool) {
	if requestInfo.Type == Requests.LoginRequest || requestInfo.Type == Requests.SignupRequest {
		if responeInfo.Type == ValidRespone {
			user := helper.GetEncodedUser(requestInfo.RequestData)
			loggedUsers := Manager.GetLoggedUsers()
			for _, loggedUser := range loggedUsers {
				if loggedUser.IsEquals(user) {

				}
			}
		}
	}
	return nil, false
}
