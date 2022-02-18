package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Zone struct {
	ID   int
	Name string
	Rooms []*Room
}

type Room struct{
	ID int
	Zone *Zone
	Name string
	Description string
	Exists [6]Exit
}

type Exit struct{
	To *Room
	Description string
}

func openDatabase(databasePath string) *sql.DB {
	path := "world.db"
	options := "?" + "_busy_timeout=10000" +
		"&" + "_foreign_keys=ON"
	db, err := sql.Open("sqlite3", path+options)
	if err != nil {
		// handle error
	}
	return db
}
