package main

import (
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *telnetClient) Connect() error {
	var err error
	if c.conn, err = net.DialTimeout("tcp", c.address, c.timeout); err != nil {
		return err
	}
	log.Printf("...Connected to %s\n", c.address)
	return nil
}

func (c *telnetClient) Close() error {
	return c.conn.Close()
}

func (c *telnetClient) Send() error {
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return err
	}
	log.Println("...EOF")
	return nil
}

func (c *telnetClient) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return err
	}
	log.Printf("...Connection was close by peer")
	return nil
}
