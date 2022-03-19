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
		go handleConnection(conn, from_player)
	}
}

func handleConnection(conn net.Conn, writeChan chan PlayerEvent) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	for {
		if !checkPlayerConnInWorld(conn) {
			//netData, err := bufio.NewReader(conn).ReadString('\n')
			//if err != nil {
			//fmt.Println(err)
			//return
			//}
			//temp := strings.TrimSpace(string(netData))
			//if temp == "STOP" {
			//break
			//}
			//fmt.Printf("Con Addr: %s\n Conn: %v\n", conn.RemoteAddr().String(), conn)
			//fmt.Printf("PLAYERS: %v\n", PLAYERS)
			player := createPlayer(conn)
			addPlayerToWorld(player)
			fmt.Printf("introducing %s to world.\n", player.Name)
			go player.captureMudEvents()
			go introducePlayerToWorld(conn, player, writeChan)
		}
	}
	fmt.Printf("closing connection: %v", conn.LocalAddr().String())
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
	player := &Player{"rantikurim", 3001, conn, conn.RemoteAddr().String(), nil}
	player.Name = getPlayerInput(conn, player, "Name? \n>")
	player.to_Player = make(chan MudEvent, 1)
	return player
}

// TODO Gets waits for user input then sends input to main go routine through shared channel
func introducePlayerToWorld(conn net.Conn, player *Player, writeChan chan PlayerEvent) {
	//cmdLook(player, "look")
	player.Printf("Adding player to world...\n")
	log.SetFlags(log.Ltime | log.Lshortfile)
	//	if err := playerCommandloop(conn, player, writeChan); err != nil {
	//		log.Fatalf("%v", err)
	//	}
	for {
		if err := playerCommandloop(conn, player, writeChan); err != nil {
			log.Fatalf("%v", err)
		}
	}

}

func getRoomString(roomId int) string {
	ret := ""
	exitsString := "[ Exits: "
	room := ROOMS[roomId]
	ret += room.Name + "\n\n"
	//?? .Description anything seems to come with a \n char?
	ret += room.Description
	for i := range room.Exits {
		// exit exists in direction i
		if room.Exits[i] != (Exit{}) {
			direction := exitIndextoDirection(i)
			exitsString += direction + " "
		}
	}
	exitsString += "]"
	ret += exitsString + "\n"
	return ret
}

func writeRoomToChannel(player *Player, roomId int) {
	s := getRoomString(roomId)
	me := MudEvent{}
	me.event = s
	player.to_Player <- me
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
