package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("Incorrect number of arguments. Host and port are required")
	}
	host := net.JoinHostPort(args[0], args[1])

	tClient := NewTelnetClient(host, timeout, os.Stdin, os.Stdout)
	err := tClient.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = tClient.Close()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := tClient.Send()
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		err := tClient.Receive()
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		sChan := make(chan os.Signal, 1)
		signal.Notify(sChan, syscall.SIGINT, syscall.SIGTERM)

		// wait for signal in channel
		<-sChan
		// call cancel()
		cancel()
	}()

	<-ctx.Done()
}
