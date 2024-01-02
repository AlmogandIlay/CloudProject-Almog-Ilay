package RequestHandlers

import (
	"CloudDrive/authentication"
)

type LoginManagerFactory struct {
	manager *authentication.LoginManager
}

func newLoginManagerFactory() (LoginManagerFactory, error) {
	manager, err := authentication.InitializeLoginManager()
	if err != nil {
		return LoginManagerFactory{}, err
	}

	return LoginManagerFactory{manager}, nil
}

func (factory LoginManagerFactory) getManager() authentication.LoginManager {
	return *factory.manager
}
