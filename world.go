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
	Exits       [6]Exit
}

type Exit struct {
	ToRoom      *Room
	Description string
}

// converts string direction of exit to Room exit index
func exitDirectionToIndex(direction string) int {
	// all Rooms have a potential of six exits
	switch direction {
	case "n":
		return 0
	case "e":
		return 1
	case "w":
		return 2
	case "s":
		return 3
	case "u":
		return 4
	case "d":
		return 5
	default:
		return -1
	}
}

func exitIndextoDirection(index int) string {
	switch index {
	case 0:
		return "n"
	case 1:
		return "e"
	case 2:
		return "w"
	case 3:
		return "s"
	case 4:
		return "u"
	case 5:
		return "d"
	default:
		return ""
	}
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

func initWorld() error {
	// Open Database and create an active transaction
	db := openDatabase(WDB)
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error while opening database: %v", err)
	}
	// Read entire world from database under one transaction
	if err = readWorld(tx, ZONES, ROOMS); err != nil {
		return fmt.Errorf("Error reading world from database: %v", err)
	}
	// Load game into memory
	tx.Commit()
	return nil
}

func readWorld(tx *sql.Tx, zones map[int]*Zone, rooms map[int]*Room) error {
	fmt.Println(getDateTime() + "Reading world file")
	var err error
	var readZones int
	var readRooms int
	// Read all zones from disk to mem
	if readZones, err = readAllZones(tx, zones); err != nil {
		return err
	}
	// Read all rooms from disk to mem
	if readRooms, err = readAllRooms(tx, rooms, zones); err != nil {
		return err
	}
	// Read all exits from disk to mem
	if err = readAllExists(tx, rooms); err != nil {
		return err
	}
	fmt.Printf(getDateTime()+"read %v zones and %v rooms\n", readZones, readRooms)
	return nil
}

// Prints room data from database with corresponding id (Each room id should be unique)
func readRoom(id int64) error {
	var roomID int
	var roomName string
	var roomDescription string
	database := openDatabase(WDB)
	rows, err := database.Query("SELECT id, name, description FROM rooms WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("Error while reading room, %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&roomID, &roomName, &roomDescription)
		if err != nil {
			return fmt.Errorf("Error while scanning room, %v", err)
		}
		r := Room{
			ID:          roomID,
			Name:        roomName,
			Description: roomDescription,
		}
		fmt.Printf("ID: %v\n\nName: %v\nDescription: %v\n", r.ID, r.Name, r.Description)
	}
	return nil
}

func readAllZones(tx *sql.Tx, m map[int]*Zone) (int, error) {
	var zoneID int
	var zoneName string
	count := 0
	rows, err := tx.Query("SELECT id, name FROM zones ORDER BY id")
	if err != nil {
		return count, fmt.Errorf("Error while querying all zones, %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&zoneID, &zoneName)
		if err != nil {
			return count, fmt.Errorf("Error while scanning zones, %v\n", err)
		}
		z := &Zone{
			ID:   zoneID,
			Name: zoneName,
		}
		m[zoneID] = z
		count++
		//fmt.Printf("ID: %v\nName: %v\n\n", z.ID, z.Name)
	}
	return count, nil
}

func readAllRooms(tx *sql.Tx, rooms map[int]*Room, zones map[int]*Zone) (int, error) {
	var roomID int
	var roomZoneID int
	var roomName string
	var roomDescription string
	count := 0
	rows, err := tx.Query("SELECT id, zone_id, name, description FROM rooms ORDER BY id")
	if err != nil {
		return count, fmt.Errorf("Error while querying all rooms: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&roomID, &roomZoneID, &roomName, &roomDescription)
		if err != nil {
			return count, fmt.Errorf("Error while scanning rooms: %v", err)
		}
		r := &Room{
			ID:          roomID,
			Zone:        zones[roomZoneID],
			Name:        roomName,
			Description: roomDescription,
		}
		rooms[roomID] = r
		count++
	}
	return count, nil
}

// Query DB for all exits and save them to corresponding rooms
func readAllExists(tx *sql.Tx, rooms map[int]*Room) error {
	var exitFromRoomID int
	var exitToRoomID int
	var exitDirection string
	var exitDescription string
	rows, err := tx.Query("SELECT from_room_id, to_room_id, direction, description FROM exits ORDER BY from_room_id")
	if err != nil {
		return fmt.Errorf("Error while querying all exits: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&exitFromRoomID, &exitToRoomID, &exitDirection, &exitDescription)
		if err != nil {
			return fmt.Errorf("Error while scanning exits: %v", err)
		}
		e := Exit{
			ToRoom:      rooms[exitToRoomID],
			Description: exitDescription,
		}
		// sets the description and destination of a single exit in a room
		rooms[exitFromRoomID].Exits[exitDirectionToIndex(exitDirection)] = e
	}
	return nil
}

func exitExists(roomId int, direction string) bool {
	room := ROOMS[roomId]
	exists := false
	for i := range room.Exits {
		if exitIndextoDirection(i) == direction {
			if room.Exits[i] != (Exit{}) {
				exists = true
			}
		}
	}
	//	if !exists {
	//		fmt.Println("You can't go that direction.")
	//	}
	return exists
}

func addPlayerToWorld(player *Player) {
	PLAYERS = append(PLAYERS[:], player)
}

func printRoom(roomId int) {
	exitsString := "[ Exits: "
	room := ROOMS[roomId]
	fmt.Println(room.Name + "\n")
	//?? .Description anything seems to come with a \n char?
	fmt.Print(room.Description)
	for i := range room.Exits {
		// exit exists in direction i
		if room.Exits[i] != (Exit{}) {
			direction := exitIndextoDirection(i)
			exitsString += direction + " "
		}
	}
	exitsString += "]"
	fmt.Println(exitsString)
}

func printExitDescription(roomId int, direction string) {
	room := ROOMS[roomId]
	exitIdx := exitDirectionToIndex(direction)
	fmt.Print(room.Exits[exitIdx].Description)
}
