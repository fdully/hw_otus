package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	// Place your code here
	currentTime := time.Now()
	exactTime, err := ntp.Time("time3.google.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("current time:", currentTime.String())
	fmt.Println("exact time:", exactTime.String())
}
