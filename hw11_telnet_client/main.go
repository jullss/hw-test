package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout for connection")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("Use: go-telnet [--timeout=10s] host port")
	}

	host, port := args[0], args[1]
	address := net.JoinHostPort(host, port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Printf("failed to connect: %v", err)
		return
	}

	defer client.Close()

	log.Printf("...Connected to %s", address)

	sendErrCh := make(chan error, 1)
	recErrCh := make(chan error, 1)

	go func() {
		sendErrCh <- client.Send()
	}()

	go func() {
		recErrCh <- client.Receive()
	}()

	select {
	case <-ctx.Done():
		fmt.Fprintln(os.Stderr, "...Stopped by signal")
	case err := <-sendErrCh:
		if err != nil {
			fmt.Fprintf(os.Stderr, "...Error: %v\n", err)
			return
		} else {
			fmt.Fprintln(os.Stderr, "...EOF")
		}
	case err := <-recErrCh:
		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Fprintf(os.Stderr, "...Error: %v\n", err)
			return
		} else {
			fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
		}
	}
}
