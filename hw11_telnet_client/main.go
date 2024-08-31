package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

type config struct {
	host    string
	port    string
	timeout time.Duration
	address string
}

func (c *config) check() error {
	if c.host == "" {
		return errors.New("host must be provided")
	}
	if c.port == "" {
		return errors.New("port must be provided")
	}
	return nil
}

func parseConfigFromOSArgs() (*config, error) {
	var cfg config
	const defaultTimeout = 10 * time.Second

	flag.StringVar(&cfg.host, "host", "", "host to connect to")
	flag.StringVar(&cfg.port, "port", "", "port to connect to")
	flag.DurationVar(&cfg.timeout, "timeout", defaultTimeout, "timeout")
	flag.Parse()

	err := cfg.check()
	if err != nil {
		return nil, err
	}
	cfg.address = net.JoinHostPort(cfg.host, cfg.port)
	return &cfg, nil
}

func main() {
	cfg, err := parseConfigFromOSArgs()
	if err != nil {
		fmt.Printf("fail to parse config: %v\n", err)
		os.Exit(-1)
	}

	client := NewTelnetClient(cfg.address, cfg.timeout, os.Stdin, os.Stdout)
	if err = client.Connect(); err != nil {
		fmt.Printf("fail to connect: %v\n", err)
		os.Exit(-1)
	}
	go func() {
		_ = client.Send()
	}()
	go func() {
		_ = client.Receive()
	}()
	signalCh := make(chan os.Signal)
	errCh := make(chan error)
	signal.Notify(signalCh, os.Interrupt, os.Kill)
	go func() {
		for sign := range signalCh {
			switch sign {
			case os.Interrupt:
				_ = client.Close()
			case os.Kill:
				_ = client.Close()
			}
		}
	}()
	<-errCh

}
