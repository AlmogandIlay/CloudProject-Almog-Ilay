package authentication

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	*sql.DB
}

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

func (db *Database) addUser(user *User) error {
	_, err := db.Exec("INSERT INTO users (userId, username, password, email) VALUES (NULL, ?, ?, ?)",
		user.username(), user.password(), user.email())
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
	var userId int
	err := db.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&userId, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *Database) getUserID(username string) (uint32, error) {
	var userId uint32
	err := db.QueryRow("SELECT userId FROM users WHERE username = ?", username).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (db *Database) getUserName(id uint32) (string, error) {
	var userName string
	err := db.QueryRow("SELECT username FROM users WHERE userId = ?", id).Scan(&userName)
	if err != nil {
		return "", err
	}
	return userName, nil
}
