package main

import (
	"bufio"
	"fmt"
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
	addCommand("quit", cmdQuit)
	addCommand("north", cmdNorth)
	addCommand("east", cmdEast)
	addCommand("west", cmdWest)
	addCommand("south", cmdSouth)
	addCommand("up", cmdUp)
	addCommand("down", cmdDown)
	addCommand("look", cmdLook)
	addCommand("recall", cmdRecall)

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
		} else {
			// playe
			if closed {
				CLOSECONNS = append(CLOSECONNS, player.Conn.RemoteAddr().String())
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
			p.writeToChannel("Huh?\n")
		}
	}
	return nil
}

// DIRECTIONS
func cmdNorth(p *Player, s string) {
	if exitExists(p.currentRoomId, "n") {
		p.doExit("n")
	} else {
		p.writeToChannel("You cannot go that way\n")
	}
}

func cmdEast(p *Player, s string) {
	if exitExists(p.currentRoomId, "e") {
		p.doExit("e")
	} else {
		p.writeToChannel("You cannot go that way\n")
	}
}

func cmdWest(p *Player, s string) {
	if exitExists(p.currentRoomId, "w") {
		p.doExit("w")
	} else {
		p.writeToChannel("You cannot go that way\n")
	}
}

func cmdSouth(p *Player, s string) {
	if exitExists(p.currentRoomId, "s") {
		p.doExit("s")
	} else {
		p.writeToChannel("You cannot go that way\n")
	}
}

func cmdUp(p *Player, s string) {
	if exitExists(p.currentRoomId, "u") {
		p.doExit("u")
	} else {
		p.writeToChannel("You cannot go that way\n")
	}
}

func cmdDown(p *Player, s string) {
	if exitExists(p.currentRoomId, "d") {
		p.doExit("d")
	} else {
		p.writeToChannel("You cannot go that way\n")
	}
}

func cmdLook(p *Player, s string) {
	words := strings.Fields(s)
	// direction to look was specified
	if len(words) > 1 {
		direction := words[1]
		writeExitDescToChannel(p, p.currentRoomId, direction)
	} else {
		p.writeRoomToChannel(p.currentRoomId)
	}
}

func cmdRecall(p *Player, s string) {
	p.doRecall()
}

// close players comunication channel, allow server go routines to terminate
// then remove player from world/data structures
func cmdQuit(p *Player, s string) {
	me := MudEvent{}
	me.event = "$"
	p.to_Player <- me
	close(p.to_Player)
	p.to_Player = nil
	// removing player from data structure now means we have active go routines for invalid player
}
