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

func playerCommandloop(player *Player, writeChan chan PlayerEvent) {
	scanner := bufio.NewScanner(player.Conn)
	// wait for player input then send to main go routine through channel
	player.Printf(">")
	for scanner.Scan() {
		//player.Printf("Scanned input: %s\n", scanner.Text())
		line := scanner.Text()
		event := PlayerEvent{}
		event.player = player
		event.Command = line
		event.Close = false

		go func() {
			writeChan <- event
		}()
	}
	if err := scanner.Err(); err != nil {
		// respond to connection being closed
		closeEvent := PlayerEvent{}
		closeEvent.player = player
		closeEvent.Command = "$"
		closeEvent.Close = true
		go func() {
			writeChan <- closeEvent
		}()
		// log that players command loop has stopped (returned)
		return
	}
}

func (p *Player) introducePlayerToWorld(writeChan chan PlayerEvent) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	cmdLook(p, "look")
	playerCommandloop(p, writeChan)
}

func (p *Player) writeToChannel(s string) {
	me := MudEvent{}
	me.event = s
	p.to_Player <- me
}

func (p *Player) writeRoomToChannel(roomId int) {
	s := getRoomString(p, roomId)
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
		// When channel is closed a "$" is written  to its queue
		me := <-player.to_Player
		if me.event != "$" {
			player.Printf(me.event)
			player.Printf(">")
		} else {
			// The players connection closes then removes the player from data structure as the player is now invalid
			fmt.Printf("Closing conn %s\n", player.Conn.RemoteAddr().String())
			player.Conn.Close()
		}
	}
}
