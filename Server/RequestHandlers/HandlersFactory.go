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

// Initialize Authentication Manager instance
func InitializeIdentifyManagerFactory() (*authentication.IdentityManager, error) {
	var err error
	Manager, err = authentication.InitializeIdentifyManager()
	if err != nil {
		Manager = &authentication.AuthenticationManager{}
		return Manager, err
	}

	return Manager, nil
}

// Access Authentication Manager instance
func GetManager() *authentication.AuthenticationManager {
	return Manager
}

func GetLoggedUser(requestInfo Requests.RequestInfo) (FileSystem.LoggedUser, error) {
	var loggedUser *FileSystem.LoggedUser
	user := helper.GetEncodedUser(requestInfo.RequestData)
	id, err := Manager.GetUserID(user.Username)
	if err != nil {
		return FileSystem.LoggedUser{}, err
	}
	loggedUser, err = FileSystem.NewLoggedUser(id)
	if err != nil {
		return FileSystem.LoggedUser{}, err
	}
	return *loggedUser, nil
}

// Remove Online User from the Online Users slice
func RemoveOnlineUser(loggedUser FileSystem.LoggedUser) error {
	username, err := Manager.GetUserName(loggedUser.UserID)
	if err != nil {
		return err
	}
	err = Manager.DeleteUser(username)
	if err != nil {
		return err
	}
	return nil
}
