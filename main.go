package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var COMMANDS = make(map[string]func(string))

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	initCommands()
	if err := commandLoop(); err != nil {
		log.Fatalf("%v", err)
	}
}

func commandLoop() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		doCommand(line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("in main command loop: %v", err)
	}
	return nil
}

func addCommand(cmd string, f func(string)) {
	for i := range cmd {
		if i == 0 {
			continue
		}
		prefix := cmd[:i]
		COMMANDS[prefix] = f
	}
	COMMANDS[cmd] = f
}

func initCommands() {
	addCommand("south", cmdSouth)
}

func doCommand(cmd string) error {
	words := strings.Fields(cmd)
	if len(words) == 0 {
		return nil
	}
	if f, exists := COMMANDS[strings.ToLower(words[0])]; exists {
		f(cmd)
	} else {
		fmt.Printf("Huh?\n")
	}
	return nil
}

func cmdSouth(s string) {
	fmt.Printf("South: %v\n", s)
}
