package world

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func openDatabase(databasePath string) *sql.DB{
	path := "world.db"
	options := "?" + "_busy_timeout=10000" +
		"&" + "_foreign_keys=ON"
	db, err := sql.Open("sqlite3", path+options)
	if err != nil {
		// handle error
	}
	return db
}
