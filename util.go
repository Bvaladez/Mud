package main

import (
	"fmt"
	"strings"
	"time"
)

func getDateTime() string {

	currentTime := time.Now()
	return fmt.Sprint(currentTime.Format("2006/01/02 15:04:05" + " "))
}

// if the string contains space as deliminiters throw out first "word"
func removeFirstWord(s string) string {
	words := strings.Fields(s)
	ss := ""
	if len(words) > 1 {
		words = append(words[1:])
		for _, word := range words {
			ss += word + " "
		}
	}
	return ss
}
