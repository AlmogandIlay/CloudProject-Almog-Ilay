package authentication

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func newDatabase() (*Database, error) {
	sqlDatabase := Database{}

	// Open the database
	db, err := sql.Open("mysql", "root:AlmogAndIlay123@tcp(127.0.0.1:3306)/cloud_db")
	if err != nil {
		return nil, err
	}
	sqlDatabase.DB = db
	return &sqlDatabase, nil
}
