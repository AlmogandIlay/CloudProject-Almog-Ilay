package main

import "fmt"

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
