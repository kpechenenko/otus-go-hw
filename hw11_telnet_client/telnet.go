package main

import (
	"bufio"
	"io"
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
	return nil
}

func (c *telnetClient) Close() error {
	return c.conn.Close()
}

func (c *telnetClient) Send() error {
	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		b := scanner.Bytes()
		if _, err := c.conn.Write(b); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (c *telnetClient) Receive() error {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		b := scanner.Bytes()
		if _, err := c.out.Write(b); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
