package main

import (
	"fmt"
	"log"
	"net"
)

type Player struct {
	Name          string
	currentRoomId int
	Conn          net.Conn
	Id	 string
}

// Makes player move from current room to next room in direction
func (p *Player) doExit(direction string) {
	toRoom := ROOMS[p.currentRoomId].Exits[exitDirectionToIndex(direction)].ToRoom
	p.currentRoomId = toRoom.ID
	PrintRoomToPlayer(p, p.currentRoomId)
}

func (p *Player) doRecall() {
	p.currentRoomId = 3001
	p.Printf("You pray to your god. Your vision blurs briefly.\n")
	PrintRoomToPlayer(p, p.currentRoomId)
}

func (p *Player) Printf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	_, err := fmt.Fprint(p.Conn, msg)
	if err != nil {
		log.Printf("network error while printing: %v", err)
	}
}
