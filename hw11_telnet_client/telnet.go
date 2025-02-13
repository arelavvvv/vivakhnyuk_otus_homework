package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient struct {
	address  string
	timeout  time.Duration
	conn     net.Conn
	stdin    io.ReadCloser
	stdout   io.Writer
	done     chan error
	err      error
	isClosed bool
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) *TelnetClient {
	return &TelnetClient{
		address:  address,
		timeout:  timeout,
		stdin:    in,
		stdout:   out,
		done:     make(chan error, 2),
		isClosed: false,
	}
}

func (c *TelnetClient) Connect() error {
	dialer := net.Dialer{Timeout: c.timeout}
	conn, err := dialer.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", c.address, err)
	}
	c.conn = conn
	log.Printf("...Connected to %s", c.address)
	return nil
}

func (c *TelnetClient) Close() error {
	if c.isClosed {
		return nil
	}
	c.isClosed = true
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return err
		}
		log.Println("...Connection was closed")
	}
	return nil
}

func (c *TelnetClient) Send() error {
	defer func() {
		c.done <- c.err
	}()

	reader := bufio.NewReader(c.stdin)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("...EOF")
				return nil
			}
			return fmt.Errorf("error reading from stdin: %w", err)
		}

		_, err = io.WriteString(c.conn, str)
		if err != nil {
			return fmt.Errorf("error writing to connection: %w", err)
		}
	}
}

func (c *TelnetClient) Receive() error {
	defer func() {
		c.done <- c.err
	}()

	_, err := io.Copy(c.stdout, c.conn)
	if err != nil {
		return err
	}

	return nil
}
