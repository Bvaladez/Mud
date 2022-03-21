//ƁĀRĐEĐ
//ƁĀƦĐƎĐ
//ƁĀƦĐƎĐ
//ƁĀƦĐΞĐ
//βѦ℞ĐΞĐ⎛
package main

var COMMANDS = make(map[string]func(*Player, string))
var ZONES = make(map[int]*Zone)
var ROOMS = make(map[int]*Room)
var PLAYERS = []*Player{}
var CLOSECONNS []string

func main() {
	// MAIN GO ROUTINE
	initCommands()
	initWorld()
	go serverServe() // TODO serverServer returns an error
	capturePlayerCommands()
}
