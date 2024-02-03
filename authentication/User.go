package authentication

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
)

// for Valid interface
type Name string
type Password string
type Email string

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// every field need to implement "Valid" method to check if the value is valid
type Validator interface {
	Valid() bool
}

// constants
const (
	MINIMUM_PASSWORD_LENGTH = 8
	MAXIMUM_PASSWORD_LENGTH = 16
	MINIMUM_USERNAME_LENGTH = 4
	MAXIMUM_USERNAME_LENGTH = 16
)

func NewUser(username, password, email string) (*User, []error) {
	var validationErrors []error

	// check if the interface implementators arent valid and append to the errors slice
	validate := func(validField Validator, err error) {
		if !validField.Valid() {
			validationErrors = append(validationErrors, err)
		}
	}

	// check for user value before creation
	validate(Name(username), &UsernameError{username})
	validate(Password(password), &PasswordError{password})
	validate(Email(email), &EmailError{email})

	// check if an error accured
	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	hashedPassword := Hash(password)

	return &User{Username: username, Password: hashedPassword, Email: email}, nil
}

// getters for each one of the User struct fields
func (user *User) username() string {
	return user.Username
}
func (user *User) password() string {
	return user.Password
}
func (user *User) email() string {
	return user.Email
}

// setters for each one of the User struct fields
func (user *User) setName(newName string) error {
	var err error = nil
	if Name(newName).Valid() {
		user.Username = newName
	} else {
		err = &UsernameError{newName}
	}
	return err
}
func (user *User) setPassword(newPassword string) error {
	var err error = nil
	if Password(newPassword).Valid() {
		user.Password = Hash(newPassword)
	} else {
		err = &PasswordError{newPassword}
	}
	return err
}
func (user *User) setEmail(newEmail string) error {
	var err error = nil
	if Email(newEmail).Valid() {
		user.Email = newEmail
	} else {
		err = &EmailError{newEmail}
	}
	return err
}

// implementation of the interface, check if the fields are valid
func (name Name) Valid() bool {
	return len(name) >= MINIMUM_USERNAME_LENGTH && len(name) <= MAXIMUM_USERNAME_LENGTH
}

func (password Password) Valid() bool {
	return len(password) >= MINIMUM_PASSWORD_LENGTH && len(password) <= MAXIMUM_PASSWORD_LENGTH
}

// check the validication of an email by email rules
func (email Email) Valid() bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, string(email))
	return match
}

// Hash sha256 encrypts the password
func Hash(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash
}

func (user *User) IsEquals(other *User) bool {
	return user.username() == other.username() && user.password() == Hash(other.password()) && user.email() == other.email()
}
