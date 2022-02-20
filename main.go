package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var COMMANDS = make(map[string]func(string))

func main() {
	var ZONES = make(map[int]*Zone)
	var ROOMS = make(map[int]*Room)
	// database not being closed anywhere potential problem?
	db := openDatabase(WDB)
	// Read all zones from disk to mem
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	if err = readAllZones(tx, ZONES); err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	// Read all rooms from disk to mem
	tx, err = db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	if err = readAllRooms(tx, ROOMS, ZONES); err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	// Read all exits from disk to mem

	log.SetFlags(log.Ltime | log.Lshortfile)
	initCommands()
	if err := commandLoop(); err != nil {
		log.Fatalf("%v", err)
	}
}

func commandLoop() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		doCommand(line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("in main command loop: %v", err)
	}
	return nil
}

func addCommand(cmd string, f func(string)) {
	for i := range cmd {
		if i == 0 {
			continue
		}
		prefix := cmd[:i]
		COMMANDS[prefix] = f
	}
	COMMANDS[cmd] = f
}

func initCommands() {
	addCommand("south", cmdSouth)
}

func doCommand(cmd string) error {
	words := strings.Fields(cmd)
	if len(words) == 0 {
		return nil
	}
	if f, exists := COMMANDS[strings.ToLower(words[0])]; exists {
		f(cmd)
	} else {
		fmt.Printf("Huh?\n")
	}
	return nil
}

func cmdSouth(s string) {
	fmt.Printf("South: %v\n", s)
}
