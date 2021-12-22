package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-hclog"
	"github.com/maticnetwork/ethstats-backend/ethstats"
)

var (
	defaultDBEndpoint = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
)

func main() {
	var wsAddr, dbEndpoint, logLevel, frontendAddr string

	flag.StringVar(&dbEndpoint, "db-endpoint", defaultDBEndpoint, "")
	flag.StringVar(&wsAddr, "ws-addr", "localhost:8000", "ws service address for collector")
	flag.StringVar(&logLevel, "log-level", "Log level", "info")
	flag.StringVar(&frontendAddr, "frontend-addr", "", "")
	flag.Parse()

	config := &ethstats.Config{
		Endpoint:      dbEndpoint,
		CollectorAddr: wsAddr,
		FrontendAddr:  frontendAddr,
	}

	logger := hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString(logLevel)})
	srv, err := ethstats.NewServer(logger, config)
	if err != nil {
		fmt.Printf("[ERROR]: %v", err)
		os.Exit(0)
	}

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	<-signalCh
	srv.Close()
}
