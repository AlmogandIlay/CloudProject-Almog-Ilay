package main

type LoginManager struct {
	*Database
	loggedUsers []User
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

func (manager *LoginManager) signin(user *User) {
	userExist, err := manager.doesUserExist(user.Username())
	if err != nil {

	}
	if userExist {

	}
	manager.addUser(user.Username(), user.Password(), user.Email())
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
