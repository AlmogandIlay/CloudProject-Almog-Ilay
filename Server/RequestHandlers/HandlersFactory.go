package RequestHandlers

import (
	"sync"
	"CloudDrive/authentication"
)

var (
	manager	*authentication.LoginManager
	Once	sync.Once
)

func InitializeFactory() (*authentication.LoginManager, error) {
	
	manager, err = authentication.InitializeLoginManager()
	if err != nil {
		manager = LoginManager{}
		return manager, err
	}

	return manager, nil
}

func GetManager() (*authentication.LoginManager, error) {
	sync.Do(func()) {
		manager, err = InitializeFactory()
		if err != nil {
			return manager, err
		}
	}
	return &manager, nil
}
