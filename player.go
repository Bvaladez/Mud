package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type MudEvent struct {
	event string
}

type Player struct {
	Name          string
	currentRoomId int
	Conn          net.Conn
	Id            string
	to_Player     chan MudEvent
}

func playerCommandloop(conn net.Conn, player *Player, writeChan chan PlayerEvent) error {
	scanner := bufio.NewScanner(conn)
	// wait for player input then send to main go routine through channel
	player.Printf(">")
	for scanner.Scan() {
		//player.Printf("Scanned input: %s\n", scanner.Text())
		line := scanner.Text()
		event := PlayerEvent{}
		event.player = player
		event.Command = line

		go func() {
			writeChan <- event
		}()

	}
	player.Printf("Stopped scanning on player command loop")
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error in player main command loop:\nE:%v\n", err)
	}
	return nil
}

func (p *Player) writeToChannel(s string) {
	me := MudEvent{}
	me.event = s
	p.to_Player <- me
}

func (p *Player) writeRoomToChannel(roomId int) {
	s := getRoomString(roomId)
	me := MudEvent{}
	me.event = s
	p.to_Player <- me
}

// Makes player move from current room to next room in direction
func (p *Player) doExit(direction string) {
	toRoom := ROOMS[p.currentRoomId].Exits[exitDirectionToIndex(direction)].ToRoom
	p.currentRoomId = toRoom.ID
	p.writeRoomToChannel(p.currentRoomId)
}

func (p *Player) doRecall() {
	p.currentRoomId = 3001
	p.writeToChannel("You pray to your god. Your vision blurs briefly.\n")
	p.writeRoomToChannel(p.currentRoomId)
}

func (p *Player) Printf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	_, err := fmt.Fprint(p.Conn, msg)
	if err != nil {
		log.Printf("network error while printing: %v", err)
	}
}

func (player *Player) captureMudEvents() {
	for {
		//p.Printf("Reading from channel\n")
		me := <-player.to_Player
		player.Printf(me.event)
		player.Printf(">")
	}
}
