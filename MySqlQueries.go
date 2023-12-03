package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	*sql.DB
}

// add methods to change values in specific row

func (db *Database) addUser(user *User) error {
	_, err := db.Exec("INSERT INTO users (id, username, password, email) VALUES (NULL, ?, ?, ?)",
		user.Username(), user.Password(), user.Email())
	return err
}

func (db *Database) removeUser(username string) error {
	_, err := db.Exec("DELETE FROM users WHERE username = ?", username)
	return err
}

func (db *Database) doesUserExist(username string) (bool, error) {
	var exist bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (db *Database) doesPasswordMatch(username string, password string) (bool, error) {
	var match bool
	err := db.QueryRow("SELECT EXISTS(SELECT * FROM users WHERE username = ? AND password = ?)", username, password).Scan(&match)
	if err != nil {
		return false, nil
	}
	return match, nil
}

func (db *Database) getUser(username string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&user.username, &user.password, &user.email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
