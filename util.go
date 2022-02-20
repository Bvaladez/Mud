package main

import (
	"fmt"
	"time"
)

func getDateTime() string {

	currentTime := time.Now()
	return fmt.Sprint(currentTime.Format("2006/01/02 15:04:05" + " "))
}
