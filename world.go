package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const WDB = "world.db"

var (
	ctx context.Context
	db  *sql.DB
)

type Zone struct {
	ID    int
	Name  string
	Rooms []*Room
}

type Room struct {
	ID          int
	Zone        *Zone
	Name        string
	Description string
	Exists      [6]Exit
}

type Exit struct {
	To          *Room
	Description string
}

func openDatabase(databasePath string) *sql.DB {
	options := "?" + "_busy_timeout=10000" +
		"&" + "_foreign_keys=ON"
	database, err := sql.Open("sqlite3", databasePath+options)
	if err != nil {
		log.Fatal(err)
	}
	return database
}

// Reads room from data base with corresponding id (Each room id should be unique)
func readRoom(id int64) {
	var roomID int
	var roomName string
	var roomDescription string
	database := openDatabase(WDB)
	rows, err := database.Query("SELECT id, name, description FROM rooms WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&roomID, &roomName, &roomDescription)
		r := Room{
			ID:          roomID,
			Name:        roomName,
			Description: roomDescription,
		}
		fmt.Printf("ID: %v\n\nName: %v\nDescription: %v\n", r.ID, r.Name, r.Description)
	}
}

func readAllZones(tx *sql.Tx, m map[int]*Zone) error {
	var zoneID int
	var zoneName string
	rows, err := tx.Query("SELECT id, name FROM zones ORDER BY id")
	if err != nil {
		return fmt.Errorf("Error while querying all zones, %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&zoneID, &zoneName)
		if err != nil {
			return fmt.Errorf("Error while scanning zones, %v\n", err)
		}
		z := &Zone{
			ID:   zoneID,
			Name: zoneName,
		}
		m[zoneID] = z
		fmt.Printf("ID: %v\nName: %v\n\n", z.ID, z.Name)
	}
	return nil
}

func readAllRooms(tx *sql.Tx, rooms map[int]*Room, zones map[int]*Zone) error {
	var roomID int
	var roomZoneID int
	var roomName string
	var roomDescription string
	rows, err := tx.Query("SELECT id, zone_id, name, description FROM rooms ORDER BY id")
	if err != nil {
		return fmt.Errorf("Error while querying all rooms: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&roomID, &roomZoneID, &roomName, &roomDescription)
		if err != nil {
			return fmt.Errorf("Error while scanning rooms: %v", err)
		}
		r := &Room{
			ID:          roomID,
			Zone:        zones[roomZoneID],
			Name:        roomName,
			Description: roomDescription,
		}
		rooms[roomID] = r
	}
	return nil
}
