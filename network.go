package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// Is controlled by a go routinne so it must handle its own errors
func serverServe() error {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		fmt.Println("Serving clients")
		conn, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("Error while waiting for client to connect: %v", err)
		}
		fmt.Println("--- NEW CONNECTION ---")
		go handleConnections(conn, from_player)
	}
}

// find new player, initialize them, capture mud events, caputre player events
func handleConnections(conn net.Conn, writeChan chan PlayerEvent) {
	player := createPlayer(conn)
	addPlayerToWorld(player)
	fmt.Printf("Welcoming %s to the world.\n", player.Name)
	go player.captureMudEvents()
	go introducePlayerToWorld(player, writeChan)
	return
}

//Create a new player object
func createPlayer(conn net.Conn) *Player {
	player := &Player{"rantikurim", 3001, conn, conn.RemoteAddr().String(), nil}
	player.Name = getPlayerInput(conn, player, "Name? \n>")
	player.to_Player = make(chan MudEvent, 3)
	return player
}

// Creates a string that can be written to player channels with default room info
func getRoomString(player *Player, roomId int) string {
	ret := ""
	exitsString := "[ Exits: "
	room := ROOMS[roomId]
	// Title name of room
	ret += room.Name + "\n\n"
	// Desciption of room
	ret += room.Description
	// Exits in room
	for i := range room.Exits {
		// exit exists in direction i
		if room.Exits[i] != (Exit{}) {
			direction := exitIndextoDirection(i)
			exitsString += direction + " "
		}
	}
	exitsString += "]"
	ret += exitsString + "\n"
	// Players in room
	ret += getPlayersString(player)
	return ret
}

// Creates a string consisting of all players in *Players current room
func getPlayersString(p *Player) string {
	playersString := "[ Players: "
	for _, storedPlayer := range PLAYERS {
		if storedPlayer.currentRoomId == p.currentRoomId {
			name := storedPlayer.Name
			playersString += name + " "
		}
	}
	playersString += "]\n"
	return playersString
}

func getExitDescString(player *Player, roomId int, direction string) string {
	room := ROOMS[roomId]
	exitIdx := exitDirectionToIndex(direction)
	if exitExists(roomId, direction) {
		return fmt.Sprintf(room.Exits[exitIdx].Description)
	} else {
		return fmt.Sprintf("You do not see anything interesting\n")
	}
}

func writeExitDescToChannel(player *Player, roomId int, direction string) {
	s := getExitDescString(player, roomId, direction)
	me := MudEvent{}
	me.event = s
	player.to_Player <- me
}

func printExitDescToPlayer(p *Player, roomId int, direction string) {
	room := ROOMS[roomId]
	exitIdx := exitDirectionToIndex(direction)
	if exitExists(roomId, direction) {
		p.Printf(room.Exits[exitIdx].Description)
	} else {
		p.Printf("You do not see anything interesting\n")
	}
}

func getPlayerInput(conn net.Conn, p *Player, s string) string {
	var input string
	scanner := bufio.NewScanner(conn)
	p.Printf(s)
	for scanner.Scan() {
		input = scanner.Text()
		fmt.Printf("Players Name: %s\n", input)
		//fmt.Printf("Player Connection: %s\n", conn.RemoteAddr().String())
		break
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error when reading users name: %v", err)
	}
	return input
}
