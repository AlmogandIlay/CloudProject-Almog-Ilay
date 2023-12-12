package main

type LoginManager struct {
	*Database
	loggedUsers []User
}

func newLoginManager() *LoginManager { // Constructor function of Login Manager
	var manager LoginManager
	var err error
	manager.loggedUsers = make([]User, 0)
	manager.Database, err = newDatabase()
	if err != nil {
		panic(err.Error())
	}
	return &manager
}

func (manager *LoginManager) login(username string, password string) {
	userExist, err := manager.doesUserExist(username)
	if err != nil {
		// query error
	}
	if !userExist {
		// user not found error
	}
	match, err := manager.doesPasswordMatch(username, password)
	if err != nil {
		// query error
	}
	if !match {
		//username does not match
	}
	user, err := manager.getUser(username)
	if err != nil {
		// query error
	}
	manager.loggedUsers = append(manager.loggedUsers, *user)
}

func (manager *LoginManager) signin(user *User) {
	userExist, err := manager.doesUserExist(user.Username())
	if err != nil {
		// query error
	}
	if userExist {
		// user already exist
	}
	manager.addUser(user.Username(), user.Password(), user.Email())
	manager.loggedUsers = append(manager.loggedUsers, *user)
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
