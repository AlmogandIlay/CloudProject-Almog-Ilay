package authentication

import (
	"strings"
)

type AuthenticationManager struct {
	*Database
	loggedUsers []User
}

const (
	ErrorFieldExistsCode = "Error 1062" // mySQL ERROR after adding exist data
)

// Constructor function of Login Manager
func InitializeAuthenticationManager() (*AuthenticationManager, error) {
	var manager AuthenticationManager
	var err error

	manager.loggedUsers = make([]User, 0)
	manager.Database, err = openDatabase()

	if err != nil {
		return nil, err
	}
	return &manager, nil
}

// Add comment
func (manager *AuthenticationManager) Login(username, password string) error {
	userExist, err := manager.doesUserExist(username)
	if err != nil {
		return err
	}
	if !userExist {
		return &UsernameNotExistsError{username}
	}

	match, err := manager.doesPasswordMatch(username, Hash(password)) // check if the password match after encryption
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

// Add comment
func (manager *AuthenticationManager) Signup(username, password, email string) []error {
	userExist, err := manager.doesUserExist(username)
	if err != nil {
		return []error{err}
	}
	if userExist {
		return []error{&UsernameExistsError{username}}
	}

	user, userErrors := NewUser(username, password, email)

	if len(userErrors) > 0 { // no errors
		return userErrors
	}

	err = manager.addUser(user) // add to the database
	if err != nil {
		if strings.Contains(err.Error(), ErrorFieldExistsCode) && strings.Contains(err.Error(), email) {
			return []error{&EmailExistsError{email}}
		}
		return []error{err}
	}
	manager.loggedUsers = append(manager.loggedUsers, *user) // add to the loggedUser slice
	return nil
}

// Add comment
func (manager *AuthenticationManager) Logout(username string) error {
	_, err := manager.doesUserExist(username)
	if err != nil {
		return err
	}
	for index, user := range manager.loggedUsers {
		if user.Username == username {
			manager.loggedUsers = append(manager.loggedUsers[:index], manager.loggedUsers[index+1:]...) // remove the user
			return nil
		}
	}
	return &UsernameNotExistsError{username}
}

// Add comment
func (manager *AuthenticationManager) DeleteUser(username string) error {
	userExist, err := manager.doesUserExist(username)
	if err != nil {
		return err
	}
	if !userExist {
		return &UsernameNotExistsError{username}
	}
	for index, user := range manager.loggedUsers { // look for the index of the user
		if user.Username == username {
			manager.loggedUsers = append(manager.loggedUsers[:index], manager.loggedUsers[index+1:]...) // remove the user
			return nil
		}
	}
	err = manager.removeUser(username)
	if err != nil {
		return err
	}
	return nil
}

func (manager *AuthenticationManager) GetLoggedUsers() []User {
	return manager.loggedUsers
}

func (manager *AuthenticationManager) GetUserID(username string) (uint32, error) {
	return manager.getUserID(username)
}
