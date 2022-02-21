package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func serverServe() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Waiting on Conn?...\nError: %v\n", err)
		}
		fmt.Println("---------------NEW CONNECTION-----------------")
		player := &Player{"BOB", 3001, conn}
		go handleConnection(conn, player)
	}
}

func handleConnection(c net.Conn, p *Player) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}
		name := getConnInput(p)
		p.Name = name
		connectPlayer(p)
		cmdLook(p, "look")
		log.SetFlags(log.Ltime | log.Lshortfile)
		if err := commandLoop(c, p); err != nil {
			log.Fatalf("%v", err)
		}

	}
	c.Close()
}

func connectPlayer(p *Player) {
	p.Printf("Connected player")
	fmt.Println("Player connecting...")
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Refused Connection?\nError: %v\n", err)
	}
	fmt.Fprint(conn, "HELLO \n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf("%s", status)
}

func getConnInput(p *Player) string {
	p.Printf("Name? ")
	return ""
}
