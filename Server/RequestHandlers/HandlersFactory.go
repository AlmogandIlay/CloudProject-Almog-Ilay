package RequestHandlers

import (
	"CloudDrive/authentication"
)

type LoginManagerFactory struct {
	manager *authentication.LoginManager
}

func NewLoginManagerFactory() (LoginManagerFactory, error) {
	manager, err := authentication.InitializeLoginManager()
	if err != nil {
		return LoginManagerFactory{}, err
	}

	return LoginManagerFactory{manager}, nil
}

func (factory LoginManagerFactory) GetManager() authentication.LoginManager {
	return *factory.manager
}
