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
		// Quit command already issued, Respond letting main mud this routine shut down
		player.writeCloseEventToMud(writeChan)
		return
	} else if !checkPlayersConnAlive(player) {
		// Disconnect, write the quit event and response to channel to queue the response behind the quit command
		player.writeCloseEventToMud(writeChan)
		player.writePlayerEventToMud(writeChan, "quit")
		return
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

func (p *Player) writeCloseEventToMud(writeChan chan PlayerEvent) {
	closeEvent := PlayerEvent{}
	closeEvent.player = p
	closeEvent.Command = "$"
	closeEvent.Close = true
	go func() {
		writeChan <- closeEvent
	}()
}

func (p *Player) writePlayerEventToMud(writeChan chan PlayerEvent, cmd string) {
	pe := PlayerEvent{}
	pe.player = p
	pe.Command = cmd
	pe.Close = false
	go func() {
		writeChan <- pe
	}()
}

func (p *Player) Printf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	_, err := fmt.Fprint(p.Conn, msg)
	if err != nil {
		log.Printf("network error while printing. Error:\n %v", err)

	}
}
