package RequestHandlers

import (
	"CloudDrive/FileSystem"
	helper "CloudDrive/Helper"
	"CloudDrive/Server/RequestHandlers/Requests"
	"CloudDrive/authentication"
)

var (
	Manager *authentication.IdentityManager
)

func InitializeIdentifyManagerFactory() (*authentication.IdentityManager, error) {
	var err error
	Manager, err = authentication.InitializeIdentifyManager()
	if err != nil {
		Manager = &authentication.IdentityManager{}
		return Manager, err
	}

	return Manager, nil
}

func GetManager() *authentication.IdentityManager {
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
