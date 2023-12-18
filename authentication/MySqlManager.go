package authentication

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	*sql.DB
}

// add methods to change values in specific row

func openDatabase() (*Database, error) {
	sqlDatabase := Database{}

	// Open the database
	db, err := sql.Open("mysql", "root:AlmogAndIlay123@tcp(127.0.0.1:3306)/cloud_db")
	if err != nil {
		return nil, err
	}
	sqlDatabase.DB = db
	return &sqlDatabase, nil
}

func (db *Database) closeDatabase() error {
	return db.Close()
}

func (db *Database) addUser(username, password, email string) error {
	_, err := db.Exec("INSERT INTO users (userId, username, password, email) VALUES (NULL, ?, ?, ?)",
		username, password, email)
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

// check if there the password match the given username
func (db *Database) doesPasswordMatch(username, password string) (bool, error) {
	var match bool
	err := db.QueryRow("SELECT EXISTS(SELECT * FROM users WHERE username = ? AND password = ?)", username, password).Scan(&match)
	if err != nil {
		return false, err
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
