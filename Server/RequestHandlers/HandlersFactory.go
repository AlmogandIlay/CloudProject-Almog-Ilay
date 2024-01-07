package RequestHandlers

import (
	"CloudDrive/authentication"
)

var (
	Manager *authentication.IdentityManager
)

func InitializeFactory() (*authentication.IdentityManager, error) {
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
