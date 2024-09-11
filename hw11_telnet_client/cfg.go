package main

import (
	"errors"
	"flag"
	"net"
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

	flag.StringVar(&cfg.host, "host", "localhost", "host to connect to")
	flag.StringVar(&cfg.port, "port", "4242", "port to connect to")
	flag.DurationVar(&cfg.timeout, "timeout", defaultTimeout, "timeout")
	flag.Parse()

	err := cfg.check()
	if err != nil {
		return nil, err
	}
	cfg.address = net.JoinHostPort(cfg.host, cfg.port)
	return &cfg, nil
}
