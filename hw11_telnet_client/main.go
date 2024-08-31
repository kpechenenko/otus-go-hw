package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	client TelnetClient
	exit   = make(chan bool, 1)
)

func main() {
	go listenSignals()
	go start()
	<-exit
}

func start() {
	var err error
	var cfg *config
	if cfg, err = parseConfigFromOSArgs(); err != nil {
		log.Fatalf("fail to parse config: %v", err)
	}
	client = NewTelnetClient(cfg.address, cfg.timeout, os.Stdin, os.Stdout)
	if err = client.Connect(); err != nil {
		log.Fatalf("fail to connect client: %v\n", err)
	}
	log.Printf("connect client to %s\n", cfg.address)
	go func() {
		if err = client.Send(); err != nil {
			log.Printf("fail to send command: %v\n", err)
		}
	}()
	go func() {
		if err = client.Receive(); err != nil {
			log.Printf("fail to receive command: %v\n", err)
		}
	}()
}

func listenSignals() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGQUIT, syscall.SIGINT)
	for {
		<-signalChan
		log.Println("stop signal received")
		err := client.Close()
		if err != nil {
			log.Fatalf("can't stop telnet client: %v", err)
		}
		exit <- true
	}
}
