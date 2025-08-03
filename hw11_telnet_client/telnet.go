package main

import (
	"errors"
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
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	connect net.Conn
}

func (c *client) Connect() error {
	var err error
	c.connect, err = net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return errors.New("connection error: " + err.Error())
	}
	log.Printf("...Connected to %s\n", c.address)
	return nil
}

func (c *client) Close() error {
	return c.connect.Close()
}

func (c *client) Send() error {
	if c.connect == nil {
		return errors.New("empty Connection")
	}

	_, err := io.Copy(c.connect, c.in)
	if errors.Is(err, io.EOF) {
		log.Println("...EOF")
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Receive() error {
	if c.connect == nil {
		return errors.New("empty Connection")
	}

	_, err := io.Copy(c.out, c.connect)
	if errors.Is(err, io.EOF) {
		log.Println("...Connection was closed by peer")
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}
