package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Place your code here
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
	var timeoutFlag string

	flag.StringVar(&timeoutFlag, "timeout", "10s", "--timeout 10s connection timeout")
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatalln("usage: go-telnet [--timeout=10s] host port")
	}

	address := net.JoinHostPort(flag.Args()[0], flag.Args()[1])

	timeout, err := time.ParseDuration(timeoutFlag)
	if err != nil {
		log.Fatalf("invalid timeout %s\n", timeoutFlag)
	}

	var osSignalCh = make(chan os.Signal, 1)
	signal.Notify(osSignalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	telnetClient := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	err = telnetClient.Connect()
	if err != nil {
		log.Fatalf("ERROR: connect %v\n", err)
	}

	go func() {
		err := telnetClient.Send()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		err := telnetClient.Receive()
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-osSignalCh

	err = telnetClient.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
