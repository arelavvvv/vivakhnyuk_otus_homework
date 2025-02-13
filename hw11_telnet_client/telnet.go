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

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}
type TelnetClientImpl struct {
	address  string
	timeout  time.Duration
	conn     net.Conn
	stdin    io.ReadCloser
	stdout   io.Writer
	done     chan error
	err      error
	isClosed bool
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientImpl{
		address:  address,
		timeout:  timeout,
		stdin:    in,
		stdout:   out,
		done:     make(chan error, 2),
		isClosed: false,
	}
}

func (c *TelnetClientImpl) Connect() error {
	dialer := net.Dialer{Timeout: c.timeout}
	conn, err := dialer.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", c.address, err)
	}
	c.conn = conn
	log.Printf("...Connected to %s", c.address)
	return nil
}

func (c *TelnetClientImpl) Close() error {
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

func (c *TelnetClientImpl) Send() error {
	reader := bufio.NewReader(c.stdin)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("...EOF")
				c.done <- nil
				return nil
			}
			c.err = fmt.Errorf("error reading from stdin: %w", err)
			c.done <- c.err
			return c.err
		}

		if _, err = io.WriteString(c.conn, str); err != nil {
			c.err = fmt.Errorf("error writing to connection: %w", err)
			c.done <- c.err
			return c.err
		}
	}
}

func (c *TelnetClientImpl) Receive() error {
	defer func() {
		c.done <- c.err
	}()

	if _, err := io.Copy(c.stdout, c.conn); err != nil {
		c.err = err
		return c.err
	}

	return nil
}
