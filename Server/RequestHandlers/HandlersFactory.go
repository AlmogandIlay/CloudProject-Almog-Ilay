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
func InitializeAuthenticationManagerFactory() error {
	var err error
	Manager, err = authentication.InitializeAuthenticationManager()
	if err != nil {
		Manager = &authentication.AuthenticationManager{}
		return err
	}

	return nil
}

// Access Authentication Manager instance
func GetManager() *authentication.AuthenticationManager {
	return Manager
}

// Initializes LoggedUser API with the given RequestInfo data
func GetLoggedUser(requestInfo Requests.RequestInfo) (FileSystem.LoggedUser, error) {
	var loggedUser *FileSystem.LoggedUser
	user := helper.GetEncodedUser(requestInfo.RequestData) // Encode RequestInfo.Data to user struct
	id, err := Manager.GetUserID(user.Username)            // Gets user id by its username
	if err != nil {                                        // If error has occured
		return FileSystem.LoggedUser{}, err
	}
	loggedUser, err = FileSystem.NewLoggedUser(id) // Initializes a new LoggedUser API Instance with the given userID
	if err != nil {                                // If error has occured
		return FileSystem.LoggedUser{}, err
	}
	return *loggedUser, nil
}

// Removes online User from the Online Users slice
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
