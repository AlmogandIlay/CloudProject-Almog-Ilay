package authentication

import "fmt"

// errors implement error interface
type UsernameError struct {
	username string
}
type PasswordError struct {
	password string
}
type EmailError struct {
	email string
}

type UsernameExistsError struct {
	username string
}
type UsernameNotExistsError struct {
	username string
}
type UsernameNotMatchPasswrodError struct {
	username string
	Password string
}

func (userError *UsernameExistsError) Error() string {
	return fmt.Sprintf("user with username '%s' already exists", userError.username)
}

func (userError *UsernameNotExistsError) Error() string {
	return fmt.Sprintf("user with username '%s' are not exists", userError.username)
}

func (userError *UsernameNotMatchPasswrodError) Error() string {
	return fmt.Sprintf("user with username '%s' not match password '%s'", userError.username, userError.Password)
}

func (userError *UsernameError) Error() string {
	return fmt.Sprintf("username '%s' is invalid! username length should be between 4-8", userError.username)
}

func (passError *PasswordError) Error() string {
	return fmt.Sprintf("password '%s' is invalid! password length should be between 8-16", passError.password)
}

func (emailError *EmailError) Error() string {
	return fmt.Sprintf("email '%s' is invalid! check for email syntex", emailError.email)
}