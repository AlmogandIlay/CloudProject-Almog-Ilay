package main

type LoginManager struct {
	*Database
}

func (manager *LoginManager) login(username string, password string) {
	userExist, err := manager.doesUserExist(username)
	if err != nil {

	}
	if !userExist {

	}
	match, err := manager.doesPasswordMatch(username, password)
	if err != nil {

	}
	if match {
		//?
	}
}

func (manager *LoginManager) signin(username string, password string, email string) {
	userExist, err := manager.doesUserExist(username)
	if err != nil {

	}
	if userExist {

	}
	manager.addUser(username, password, email)
}

func (manager *LoginManager) logout(username string) {
	userExist, err := manager.doesUserExist(username)
	if err != nil {

	}
	if !userExist {

	}
	//?
}

func (manager *LoginManager) logoutSystem(username string) {
	userExist, err := manager.doesUserExist(username)
	if err != nil {

	}
	if !userExist {

	}
	manager.removeUser(username)
}
