package main

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Name string
type Password string
type Email string

type User struct {
	username string
	password string
	email    string
}

// every field need to implement "Valid" method to check if the value is valid
type Validator interface {
	Valid() bool
}

func NewUser(username, password, email string) (*User, []error) {
	var validationErrors []error

	validate := func(validField Validator, errMsg string) {
		if !validField.Valid() {
			validationErrors = append(validationErrors, errors.New(errMsg))
		}
	}

	validate(Name(username), "invalid username")
	validate(Password(password), "invalid password")
	validate(Email(email), "invalid email")

	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	hashedPassword, err := Hash(password)
	if err != nil {
		return nil, []error{err}
	}

	return &User{username: username, password: hashedPassword, email: email}, nil
}

// getters for each one of the User struct fields
func (user *User) Username() string {
	return user.username
}
func (user *User) Password() string {
	return user.password
}
func (user *User) Email() string {
	return user.email
}

// setters for each one of the User struct fields
func (user *User) setName(newName string) error {
	var err error = nil
	if Name(user.username).Valid() {
		user.username = newName
	} else {
		err = errors.New("invalid username")
	}
	return err
}
func (user *User) setPassword(password string) error {
	var err error = nil
	if Password(user.password).Valid() {
		user.password, err = Hash(password)
	} else {
		err = errors.New("invalid password")
	}
	return err
}
func (user *User) setEmail(newEmail string) error {
	var err error = nil
	if Email(user.email).Valid() {
		user.email = newEmail
	} else {
		err = errors.New("invalid email")
	}
	return err
}

// implementation of the interface
func (name Name) Valid() bool {
	return true
}
func (password Password) Valid() bool {
	return true
}
func (email Email) Valid() bool {
	return true
}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
