package main

import (
	"flag"
	"time"

	"github.com/FeshLig/metrcollector/internal/net"
)

type Options struct {
	address        net.NetAddress
	reportInterval time.Duration
	pollInterval   time.Duration
}

func ParseFlags() Options {
	var options Options

	options.address = net.NetAddress{
		Host: "localhost",
		Port: 8080,
	}
	options.reportInterval = time.Duration(10)
	options.pollInterval = time.Duration(2)

	flag.Var(&options.address, "a", "Net address host:port")
	flag.Duration("r", options.reportInterval, "frequency of sending metrics")
	flag.Duration("p", options.pollInterval, "metrics polling frequency")

	flag.Parse()
	return options
}
