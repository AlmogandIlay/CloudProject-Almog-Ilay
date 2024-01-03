package RequestHandlers

import (
	"CloudDrive/authentication"
	"sync"
)

var (
	Manager *authentication.LoginManager
	Once    sync.Once
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
