package main

import (
	"errors"
)

type Name string
type Password string
type Email string

type User struct {
	username Name
	password Password
	email    Email
}

// every field need to implement "isValid" method to check if the value is valid
type Validatable interface {
	isValid() bool
}

func newUser(username Name, password Password, email Email) (*User, error) {

	var err error = nil
	//validate the values before create struct
	if !username.isValid() {
		err = errors.New("the username is invalid")
	}
	if !password.isValid() {
		err = errors.New("the password is invalid")
	}
	if !email.isValid() {
		err = errors.New("the email is invalid")
	}

	user := &User{username: username, password: password, email: email}

	return user, err
}

// getters for each one of the User struct fields
func (user *User) Username() Name {
	return user.username
}
func (user *User) Password() Password {
	return user.password
}
func (user *User) Email() Email {
	return user.email
}

// setters for each one of the User struct fields
func (user *User) setName(newName Name) error {
	var err error = nil
	if user.username.isValid() {
		user.username = newName
	} else {
		err = errors.New("invalid username")
	}
	return err
}
func (user *User) setPassword(newPassword Password) error {
	var err error = nil
	if user.password.isValid() {
		user.password = newPassword
	} else {
		err = errors.New("invalid password")
	}
	return err

}
func (user *User) setEmail(newEmail Email) error {
	var err error = nil
	if user.password.isValid() {
		user.email = newEmail
	} else {
		err = errors.New("invalid email")
	}
	return err
}

// implementation of the interface
func (name *Name) isValid() bool {
	return true
}
func (password *Password) isValid() bool {
	return true
}
func (email *Email) isValid() bool {
	return true
}

func (password *Password) Hash() string {
	return ""
}
