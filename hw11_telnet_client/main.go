package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func parseCommandLineArgs() (string, time.Duration, error) {
	var timeoutStr string
	flag.StringVar(&timeoutStr, "timeout", "10s", "Connection timeout")
	flag.Parse()

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid timeout format: %w", err)
	}

	args := flag.Args()
	if len(args) != 2 {
		return "", 0, fmt.Errorf("usage: go-telnet [--timeout=10s] host port")
	}

	return net.JoinHostPort(args[0], args[1]), timeout, nil
}

func setupSignalHandling(client TelnetClient) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("received SIGINT, exiting")
		client.Close()
		os.Exit(0)
	}()
}

func handleReceiveError(client TelnetClient, err error) {
	if err == nil {
		return
	}

	if !strings.Contains(err.Error(), "use of closed network connection") && !strings.Contains(err.Error(), "EOF") {
		log.Printf("error receiving: %v", err)
	}

	if closeErr := client.Close(); closeErr != nil {
		log.Printf("error closing client: %v", closeErr)
	}

	if !strings.Contains(err.Error(), "EOF") {
		log.Println("...Connection was closed by peer")
	}
}

func receiveData(client TelnetClient) {
	err := client.Receive()
	handleReceiveError(client, err)
}

func handleSendError(client TelnetClient, err error) {
	if err == nil {
		return
	}

	if !strings.Contains(err.Error(), "use of closed network connection") && !strings.Contains(err.Error(), "EOF") {
		log.Printf("error sending: %v", err)
	}

	if closeErr := client.Close(); closeErr != nil {
		log.Printf("error closing client: %v", closeErr)
	}
}

func sendData(client TelnetClient) {
	if err := client.Send(); err != nil {
		handleSendError(client, err)
	}
}

func startDataTransfer(client TelnetClient) {
	go sendData(client)
	go receiveData(client)
}

func main() {
	addr, timeout, err := parseCommandLineArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	client := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)

	err = client.Connect()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	setupSignalHandling(client)
	startDataTransfer(client)

	client.Receive()
	client.Close()
}
