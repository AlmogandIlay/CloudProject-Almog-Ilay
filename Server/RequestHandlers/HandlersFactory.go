package RequestHandlers

import (
	"CloudDrive/FileSystem"
	"CloudDrive/authentication"
)

var (
	Manager *authentication.IdentityManager
)

type CurrentLoggedUser struct {
	currentLoggedUser *FileSystem.LoggedUser
}

func InitializeIdentifyManagerFactory() (*authentication.IdentityManager, error) {
	Manager, err := authentication.InitializeIdentifyManager()
	if err != nil {
		Manager = &authentication.IdentityManager{}
		return Manager, err
	}

	return Manager, nil
}

func GetManager() *authentication.IdentityManager {
	return Manager
}

func CreateCurrentLoggedUser(id uint32) (*CurrentLoggedUser, error) {
	loggedUser, err := FileSystem.NewLoggedUser(id)
	if err != nil {
		return nil, err
	}
	return &CurrentLoggedUser{currentLoggedUser: loggedUser}, nil

}
