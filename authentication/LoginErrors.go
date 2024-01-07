package authentication

import "fmt"

// errors implement error interface
type UsernameError struct{ username string }

type PasswordError struct{ password string }

type EmailError struct{ email string }

type UsernameExistsError struct{ username string }

type UsernameNotExistsError struct{ username string }

type EmailExistsError struct{ email string }

type UsernameNotMatchPasswrodError struct {
	username string
	Password string
}

func (userError *UsernameExistsError) Error() string {
	return fmt.Sprintf("User with username '%s' already exists", userError.username)
}

func (userError *UsernameNotExistsError) Error() string {
	return fmt.Sprintf("User with username '%s' are not exists", userError.username)
}

func (userError *UsernameNotMatchPasswrodError) Error() string {
	return fmt.Sprintf("User with username '%s' not match password '%s'", userError.username, userError.Password)
}

func (userError *UsernameError) Error() string {
	return fmt.Sprintf("Username '%s' is invalid! username length should be between 4-8", userError.username)
}

func (passError *PasswordError) Error() string {
	return fmt.Sprintf("Password '%s' is invalid! password length should be between 8-16", passError.password)
}

func (emailError *EmailError) Error() string {
	return fmt.Sprintf("Email '%s' is invalid! check for email syntex", emailError.email)
}

func (userError *EmailExistsError) Error() string {
	return fmt.Sprintf("User with email '%s' already exists", userError.email)
}
