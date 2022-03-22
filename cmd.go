package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var from_player = make(chan PlayerEvent, 10)

type PlayerEvent struct {
	player  *Player
	Command string
	Close   bool
}

func initCommands() {
	// Commands prefixes get over written in the order they are added (Last is top priority)
	fmt.Println((getDateTime() + "Installing commands"))
	addCommand("say", cmdSay)
	addCommand("shout", cmdShout)
	addCommand("north", cmdNorth)
	addCommand("east", cmdEast)
	addCommand("west", cmdWest)
	addCommand("south", cmdSouth)
	addCommand("up", cmdUp)
	addCommand("down", cmdDown)
	addCommand("look", cmdLook)
	addCommand("recall", cmdRecall)
	addCommand("quit", cmdQuit)

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

func capturePlayerCommands() {
	for {
		playerEvent := <-from_player
		player := playerEvent.player
		cmd := playerEvent.Command
		closed := playerEvent.Close
		fmt.Printf("Read: %s from player: %s\n", cmd, player.Name)
		// if a player channel is nil are disconnected or disconnecting either way ignore all commands
		if player.to_Player != nil {
			err := doCommand(player, cmd)
			if err != nil {
				fmt.Printf("Error while doing player command\nCmd: %s", cmd)
			}
			// ignore what player says and remove them from data structures
			// at this point the user couldnt have typed the cmd it must be a $ signalling the player is invalid
		} else {
			if closed {
				for i, storedPlayer := range PLAYERS {
					if storedPlayer.Id == player.Id {
						PLAYERS = append(PLAYERS[:i], PLAYERS[i+1:]...)
					}
				}
			}
			continue
		}
	}
}

func commandLoop(c net.Conn, player *Player) error {
	player.Printf("Entering command loop %v\n", player.Name)
	scanner := bufio.NewScanner(c)
	player.Printf(">")
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Command: %v\n", line)
		doCommand(player, line)
		player.Printf(">")
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error in main command loop:\n E:%v\n P:%v\n", err, &player)
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

func doCommand(p *Player, cmd string) error {
	if p.Conn != nil {
		words := strings.Fields(cmd)
		if len(words) == 0 {
			return nil
		}
		if f, exists := COMMANDS[strings.ToLower(words[0])]; exists {
			f(p, cmd)
		} else {
			writeToChannel(p, "Huh?\n")
		}
	}
	return nil
}

func introducePlayerToWorld(p *Player, writeChan chan PlayerEvent) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	cmdLook(p, "look")
	playerCommandloop(p, writeChan)
}

func writeToChannel(p *Player, s string) {
	me := MudEvent{}
	me.event = s
	p.to_Player <- me
}

func writeRoomToChannel(p *Player, roomId int) {
	s := getRoomString(p, roomId)
	me := MudEvent{}
	me.event = s
	p.to_Player <- me
}

// Makes player move from current room to next room in direction
func doExit(p *Player, direction string) {
	toRoom := ROOMS[p.currentRoomId].Exits[exitDirectionToIndex(direction)].ToRoom
	p.currentRoomId = toRoom.ID
	writeRoomToChannel(p, p.currentRoomId)
}

func doRecall(p *Player) {
	p.currentRoomId = 3001
	writeToChannel(p, "You pray to your god. Your vision blurs briefly.\n")
	writeRoomToChannel(p, p.currentRoomId)
}

// DIRECTIONS
func cmdNorth(p *Player, s string) {
	if exitExists(p.currentRoomId, "n") {
		doExit(p, "n")
	} else {
		writeToChannel(p, "You cannot go that way\n")
	}
}

func cmdEast(p *Player, s string) {
	if exitExists(p.currentRoomId, "e") {
		doExit(p, "e")
	} else {
		writeToChannel(p, "You cannot go that way\n")
	}
}

func cmdWest(p *Player, s string) {
	if exitExists(p.currentRoomId, "w") {
		doExit(p, "w")
	} else {
		writeToChannel(p, "You cannot go that way\n")
	}
}

func cmdSouth(p *Player, s string) {
	if exitExists(p.currentRoomId, "s") {
		doExit(p, "s")
	} else {
		writeToChannel(p, "You cannot go that way\n")
	}
}

func cmdUp(p *Player, s string) {
	if exitExists(p.currentRoomId, "u") {
		doExit(p, "u")
	} else {
		writeToChannel(p, "You cannot go that way\n")
	}
}

func cmdDown(p *Player, s string) {
	if exitExists(p.currentRoomId, "d") {
		doExit(p, "d")
	} else {
		writeToChannel(p, "You cannot go that way\n")
	}
}

func cmdLook(p *Player, s string) {
	words := strings.Fields(s)
	// direction to look was specified
	if len(words) > 1 {
		direction := words[1]
		writeExitDescToChannel(p, p.currentRoomId, direction)
	} else {
		writeRoomToChannel(p, p.currentRoomId)
	}
}

func cmdRecall(p *Player, s string) {
	doRecall(p)
}

func cmdShout(p *Player, s string) {
	words := strings.Fields(s)
	ss := ""
	if len(words) > 1 {
		words = append(words[1:])
		for _, word := range words {
			ss += word + " "
		}
	} else {
		writestring := fmt.Sprintf("you have the thought to shout but dont do or say anything\n")
		writeToChannel(p, writestring)
		return
	}
	sUpper := strings.ToUpper(ss)
	for _, storedPlayer := range PLAYERS {
		if storedPlayer.Id != p.Id && storedPlayer.to_Player != nil {
			writeString := fmt.Sprintf("\n%s: %s\n", p.Name, sUpper)
			writeToChannel(storedPlayer, writeString)
		} else {
			writestring := fmt.Sprintf("you shouted %s\n", sUpper)
			writeToChannel(p, writestring)
		}
	}
}

// write message s to channel of all player with same current room ID.
func cmdSay(p *Player, s string) {
	s = removeFirstWord(s)
	listeningPlayers := 0
	if s != "" {
		for _, storedPlayer := range PLAYERS {
			if storedPlayer.currentRoomId == p.currentRoomId && storedPlayer.Id != p.Id {
				listeningPlayers++
				writeString := fmt.Sprintf("\n%s: %s\n", p.Name, s)
				writeToChannel(storedPlayer, writeString)
			}
		}
		writeString := fmt.Sprintf("you said %s and %d people heard\n", s, listeningPlayers)
		writeToChannel(p, writeString)
	} else {
		writestring := fmt.Sprintf("you have the thought to shout but dont do or say anything\n")
		writeToChannel(p, writestring)
	}
}

// close players comunication channel, allow server go routines to terminate
// then remove player from world/data structures
func cmdQuit(p *Player, s string) {
	me := MudEvent{}
	me.event = "$"
	p.to_Player <- me
	close(p.to_Player)
	p.to_Player = nil
}
