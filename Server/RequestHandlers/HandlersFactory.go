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

func GetLoggedUser(requestInfo Requests.RequestInfo) (FileSystem.LoggedUser, error) {
	var loggedUser *FileSystem.LoggedUser
	user := helper.GetEncodedUser(requestInfo.RequestData)
	connectedUsers := Manager.GetLoggedUsers()
	for _, connectedUser := range connectedUsers {
		if connectedUser.IsEquals(&user) {
			id, err := Manager.GetUserID(user.Username)
			if err != nil {
				return FileSystem.LoggedUser{}, err
			}
			loggedUser, err = FileSystem.NewLoggedUser(id)
			if err != nil {
				return FileSystem.LoggedUser{}, err
			}

		}
	}
	return *loggedUser, nil
}
