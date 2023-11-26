package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	*sql.DB
}

func (db *Database) addUser(user User) error {
	_, err := db.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", user.username, user.password, user.email)
	return err
}

func (db *Database) removeUser(user User) error {
	_, err := db.Exec("DELETE FROM users WHERE username = ?", user.username)
	return err
}

func (db *Database) doesUserExist(username Name) (bool, error) {
	var resultUsername string
	err := db.QueryRow("SELECT username FROM USERS WHERE username = ? LIMIT 1", username).Scan(&resultUsername)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}
