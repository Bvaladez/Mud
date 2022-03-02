package main

import (
	"fmt"
	"log"
)

var COMMANDS = make(map[string]func(*Player, string))
var ZONES = make(map[int]*Zone)
var ROOMS = make(map[int]*Room)
var PLAYERS = []*Player{}

func main() {
	// MAIN GO ROUTINE
	initCommands()
	if err := initWorld(); err != nil {
		log.Fatal(err)
	}
	// Start serving (SHOULD BE ITS OWN GO ROUTINE)
	// TODO serverServer returns an bubbled up errors and should be checked
	err := serverServe()
	// TODO Read incoming commands from network
	if err != nil {
		fmt.Print(err)
	}
}
