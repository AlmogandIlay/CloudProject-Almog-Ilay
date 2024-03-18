package authentication

import (
	"strings"
)

type AuthenticationManager struct {
	*Database
	onlineUsers []User
}

const (
	ErrorFieldExistsCode = "Error 1062" // mySQL ERROR after adding exist data
)

// Constructor function of Login Manager
func InitializeAuthenticationManager() (*AuthenticationManager, error) {
	var manager AuthenticationManager
	var err error

	manager.onlineUsers = make([]User, 0)
	manager.Database, err = openDatabase()

	if err != nil {
		return nil, err
	}
	return &manager, nil
}

// Performs login request using the AuthenticationManager API
func (manager *AuthenticationManager) Login(username, password string) error {
	userExist, err := manager.doesUserExist(username) // Checks if the username exists in the database table
	if err != nil {                                   // If error occured in the check
		return err
	}
	if !userExist { // If username does not exists
		return &UsernameNotExistsError{username}
	}

	match, err := manager.doesPasswordMatch(username, Hash(password)) // check if the password match after encryption
	if err != nil {                                                   // If error occured in the check
		return err
	}
	if !match { // If password not match
		return &UsernameNotMatchPasswrodError{username, password}
	}
	// User has succesfully signed in
	user, err := manager.getUser(username)
	if err != nil {
		return err
	}
	manager.onlineUsers = append(manager.onlineUsers, *user) // Appened the user to the onlineUsers slice
	return nil
}

// Performs signup request using the AuthenticationManager API
func (manager *AuthenticationManager) Signup(username, password, email string) []error {
	userExist, err := manager.doesUserExist(username) // Checks if the username is already exists in the database table
	if err != nil {                                   // If error occured in the check
		return []error{err}
	}
	if userExist { // If username is already exists
		return []error{&UsernameExistsError{username}}
	}

	user, userErrors := NewUser(username, password, email) // Creates the new user and validates its fields

	if len(userErrors) > 0 { // If there are errors after validating the user
		return userErrors
	}

	err = manager.addUser(user) // add to the database
	if err != nil {
		if strings.Contains(err.Error(), ErrorFieldExistsCode) && strings.Contains(err.Error(), email) {
			return []error{&EmailExistsError{email}}
		}
		return []error{err}
	}
	manager.onlineUsers = append(manager.onlineUsers, *user) // add to the onlineUsers slice
	return nil
}

// Add comment
func (manager *AuthenticationManager) Logout(username string) error {
	_, err := manager.doesUserExist(username)
	if err != nil {
		return err
	}
	for index, user := range manager.onlineUsers {
		if user.Username == username {
			manager.onlineUsers = append(manager.onlineUsers[:index], manager.onlineUsers[index+1:]...) // remove the user
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
	for index, user := range manager.onlineUsers { // looking for the username in the slices of users
		if user.Username == username {
			manager.onlineUsers = append(manager.onlineUsers[:index], manager.onlineUsers[index+1:]...) // remove the user
			return nil
		}
	}
	err = manager.removeUser(username)
	if err != nil {
		return err
	}
	return nil
}

func (manager *AuthenticationManager) GetOnlineUsers() []User {
	return manager.onlineUsers
}

func (manager *AuthenticationManager) GetUserID(username string) (uint32, error) {
	return manager.getUserID(username)
}

func (manager *AuthenticationManager) GetUserName(id uint32) (string, error) {
	return manager.getUserName(id)
}
