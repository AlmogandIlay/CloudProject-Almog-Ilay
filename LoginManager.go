package main

type LoginManager struct {
	*Database
	loggedUsers []User
}

// Constructor function of Login Manager
func newLoginManager() (*LoginManager, error) {
	var manager LoginManager
	var err error

	manager.loggedUsers = make([]User, 0)
	manager.Database, err = newDatabase()

	if err != nil {
		return nil, err
	}
	return &manager, nil
}

func (manager *LoginManager) login(username, password string) error {
	userExist, err := manager.doesUserExist(username)
	if err != nil {
		return err
	}
	if !userExist {
		return &UsernameNotExistsError{username}
	}

	match, err := manager.doesPasswordMatch(username, password)
	if err != nil {
		return err
	}
	if !match {
		return &UsernameNotMatchPasswrodError{username, password}
	}
	user, err := manager.getUser(username)
	if err != nil {
		return err
	}
	manager.loggedUsers = append(manager.loggedUsers, *user)
	return nil
}

func (manager *LoginManager) signin(username, password, email string) []error {
	userExist, err := manager.doesUserExist(username)
	if err != nil {
		return []error{err}
	}
	if userExist {
		return []error{&UsernameExistsError{username}}
	}

	user, userErrors := NewUser(username, password, email)
	for _, err := range userErrors {
		if err != nil {
			userErrors = append(userErrors, err)
		}
	}
	if len(userErrors) > 0 {
		return userErrors
	}

	manager.addUser(user.Username(), user.Password(), user.Email())
	manager.loggedUsers = append(manager.loggedUsers, *user)
	return nil
}

func (manager *LoginManager) logout(username string) error {
	_, err := manager.doesUserExist(username)
	if err != nil {
		return err
	}
	for index, user := range manager.loggedUsers {
		if user.username == username {
			manager.loggedUsers = append(manager.loggedUsers[:index], manager.loggedUsers[index+1:]...)
			return nil
		}
	}
	return &UsernameNotExistsError{username}
}

func (manager *LoginManager) deleteUser(username string) error {
	userExist, err := manager.doesUserExist(username)
	if err != nil {
		return err
	}
	if !userExist {
		return &UsernameNotExistsError{username}
	}
	manager.removeUser(username)
	return nil
}

func (manager *LoginManager) getLoggedUsers() []User {
	return manager.loggedUsers
}
