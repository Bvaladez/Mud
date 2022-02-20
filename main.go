package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var COMMANDS = make(map[string]func(*Player, string))

func main() {
	// Initialize game data
	player := &Player{"BOB", 3001}
	var ZONES = make(map[int]*Zone)
	var ROOMS = make(map[int]*Room)
	initCommands()
	// database not being closed anywhere potential problem?
	// Open Database and create an active transaction
	db := openDatabase(WDB)
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// Read entire world from database under one transaction
	if err = readWorld(tx, ZONES, ROOMS); err != nil {
		log.Fatal(err)
	}
	tx.Commit()

	log.SetFlags(log.Ltime | log.Lshortfile)
	if err := commandLoop(player); err != nil {
		log.Fatalf("%v", err)
	}
}

func commandLoop(player *Player) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		doCommand(player, line)
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
	fmt.Printf("North: %v\n", s)
}
func cmdEast(p *Player, s string) {
	fmt.Printf("East: %v\n", s)
}
func cmdWest(p *Player, s string) {
	fmt.Printf("West: %v\n", s)
}
func cmdSouth(p *Player, s string) {
	fmt.Printf("South: %v\n", s)
}
func cmdUp(p *Player, s string) {
	fmt.Printf("Up: %v\n", s)
}
func cmdDown(p *Player, s string) {
	fmt.Printf("Down: %v\n", s)
}

// INTERACTION
func cmdLook(p *Player, s string) {
	words := strings.Fields(s)
	// direction to look was specified
	if len(words) > 1 {
		direction := words[1]
		fmt.Printf("Look %v: %v\n", direction, s)
	} else {
		fmt.Printf("Look: %v\n", s)
	}
}

//ACTION
func cmdRecall(p *Player, s string) {
	fmt.Printf("Recall: %v\n", s)
}
