package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var COMMANDS = make(map[string]func(*Player, string))
var ZONES = make(map[int]*Zone)
var ROOMS = make(map[int]*Room)

func main() {
	initCommands()
	if err := initWorld(); err != nil {
		log.Fatal(err)
	}
	serverServe()
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

func getInput(s string) string {
	fmt.Println(s)
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input = scanner.Text()
	}
	return input
}


func commandLoop(c net.Conn, player *Player) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">")
	for scanner.Scan() {
		line := scanner.Text()
		doCommand(player, line)
		fmt.Print(">")
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("in main command loop: %v", err)
	}
	return nil
}

func addCommand(cmd string, f func(*Player, string)) {
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
	// Commands prefixes get over written in the order they are added (Last is top priority)
	fmt.Println((getDateTime() + "Installing commands"))
	addCommand("north", cmdNorth)
	addCommand("east", cmdEast)
	addCommand("west", cmdWest)
	addCommand("south", cmdSouth)
	addCommand("up", cmdUp)
	addCommand("down", cmdDown)
	addCommand("look", cmdLook)
	addCommand("recall", cmdRecall)

}

func doCommand(player *Player, cmd string) error {
	words := strings.Fields(cmd)
	if len(words) == 0 {
		return nil
	}
	if f, exists := COMMANDS[strings.ToLower(words[0])]; exists {
		f(player, cmd)
	} else {
		fmt.Printf("Huh?\n")
	}
	return nil
}

// DIRECTIONS
func cmdNorth(p *Player, s string) {
	if exitExists(p.currentRoomId, "n") {
		p.doExit("n")
	}
}
func cmdEast(p *Player, s string) {
	if exitExists(p.currentRoomId, "e") {
		p.doExit("e")
	}
}
func cmdWest(p *Player, s string) {
	if exitExists(p.currentRoomId, "w") {
		p.doExit("w")
	}
}
func cmdSouth(p *Player, s string) {
	if exitExists(p.currentRoomId, "s") {
		p.doExit("s")
	}
}
func cmdUp(p *Player, s string) {
	if exitExists(p.currentRoomId, "u") {
		p.doExit("u")
	}
}
func cmdDown(p *Player, s string) {
	if exitExists(p.currentRoomId, "d") {
		p.doExit("d")
	}
}

// INTERACTION
func cmdLook(p *Player, s string) {
	words := strings.Fields(s)
	// direction to look was specified
	if len(words) > 1 {
		direction := words[1]
		printExitDescription(p.currentRoomId, direction)
	} else {
		printRoom(p.currentRoomId)
	}
}

//ACTION
func cmdRecall(p *Player, s string) {
	p.doRecall()
}
