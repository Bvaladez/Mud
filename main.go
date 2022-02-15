package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0{
		}else if len(fields) > 1 {
			fmt.Printf("first word: %q, rest: %v", fields[0], fields[1:])
		}else{
			fmt.Printf("first word: %q", fields[0])
		}
	if err := scanner.Err(); err != nil {
		log.Printf("scanner err: %v", err)
	}

}
