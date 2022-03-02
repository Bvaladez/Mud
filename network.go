package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func serverServe() error {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		fmt.Println("Servering clients")
		conn, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("Error while waiting for client to connect: %v", err)
		}
		fmt.Println("--- NEW CONNECTION ---")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}
		player := createPlayer(conn)
		addPlayerToWorld(player)
		go IntroducePlayerToWorld(conn, player)
	}
	conn.Close()
}

func checkPlayerConnInWorld(conn net.Conn) bool {
	for _, player := range PLAYERS {
		if player.Conn.RemoteAddr().String() == conn.RemoteAddr().String() {
			return true
		} else {
			return false
		}
	}
	return false
}

func createPlayer(conn net.Conn) *Player {
	player := &Player{"rantikurim", 3001, conn, conn.RemoteAddr().String()}
	player.Name = getPlayerInput(conn, player, "Name? ")
	return player
}

func IntroducePlayerToWorld(conn net.Conn, player *Player) {
	cmdLook(player, "look")
	log.SetFlags(log.Ltime | log.Lshortfile)
	if err := commandLoop(conn, player); err != nil {
		log.Fatalf("%v", err)
	}
}

func PrintRoomToPlayer(p *Player, roomId int) {
	exitsString := "[ Exits: "
	room := ROOMS[roomId]
	p.Printf(room.Name + "\n\n")
	//?? .Description anything seems to come with a \n char?
	p.Printf(room.Description)
	for i := range room.Exits {
		// exit exists in direction i
		if room.Exits[i] != (Exit{}) {
			direction := exitIndextoDirection(i)
			exitsString += direction + " "
		}
	}
	exitsString += "]"
	p.Printf(exitsString + "\n")
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
		fmt.Printf("Player Connection: %s\n", conn.RemoteAddr().String())
		break
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error when reading users name: %v", err)
	}
	return input
}
