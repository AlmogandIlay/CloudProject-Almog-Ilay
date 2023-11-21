package main

import "errors"

type Name string
type Passsword string
type Email string

type User struct {
	username Name
	password Passsword
	email    Email
}

// every field need to implement "isValid" method to check if the value is valid
type Validatable interface {
	isValid() bool
}

func newUser(username Name, password Passsword, email Email) (*User, error) {

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
func (user *User) Password() Passsword {
	return user.password
}
func (user *User) Email() Email {
	return user.email
}

// setters for each one of the User struct fields
func (user *User) setName(newName Name) {
	user.username = newName
}
func (user *User) setPassword(newPassword Passsword) {
	user.password = newPassword
}
func (user *User) setEmail(newEmail Email) {
	user.email = newEmail
}

// implementation of the interface
func (name *Name) isValid() bool {
	return true
}
func (name *Passsword) isValid() bool {
	return true
}
func (name *Email) isValid() bool {
	return true
}
