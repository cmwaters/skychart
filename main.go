package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/cmwaters/skymap/server"
)

const (
	defaultUpdateFreq = "@daily"
)

func main() {
	registryUrl, listenAddr, err := parseArgs()
	if err != nil {
		fmt.Print(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	err = server.Serve(ctx, registryUrl, listenAddr, defaultUpdateFreq)
	if err != nil {
		fmt.Print(err)
	}
}

func parseArgs() (string, string, error) {
	if len(os.Args) > 3 || len(os.Args) == 1 {
		return "", "", errors.New("expected 1 or 2 arguments. \n\nUsage: skymap registry-url [listen-addr]")
	}
	registryUrl := os.Args[1]
	_, err := url.Parse(registryUrl)
	if err != nil {
		return "", "", fmt.Errorf("unable to parse registry url: %w. \n\nUsage: skymap registry-url [listen-addr]", err)
	}

	if len(os.Args) == 2 {
		return registryUrl, "", nil
	}

	return registryUrl, os.Args[2], nil
}
