package RequestHandlers

import (
	"CloudDrive/authentication"
)

var (
	Manager *authentication.LoginManager
)

func InitializeFactory() (*authentication.LoginManager, error) {
	var err error
	Manager, err = authentication.InitializeLoginManager()
	if err != nil {
		Manager = &authentication.LoginManager{}
		return Manager, err
	}

	return Manager, nil
}

func GetManager() (*authentication.LoginManager, error) {
	return Manager, nil
}
