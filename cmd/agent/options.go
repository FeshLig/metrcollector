package main

import (
	"flag"
	"time"

	net "github.com/FeshLig/metrcollector/internal/flags"
)

type Options struct {
	address        net.NetAddress
	reportInterval net.SecondsDuration
	pollInterval   net.SecondsDuration
}

func ParseFlags() Options {
	var options Options

	options.address = net.NetAddress{
		Host: "localhost",
		Port: 8080,
	}
	options.reportInterval.Duration = time.Duration(10) * time.Second
	options.pollInterval.Duration = time.Duration(2) * time.Second

	flag.Var(&options.address, "a", "net address host:port")
	flag.Var(&options.reportInterval, "r", "frequency of sending metrics")
	flag.Var(&options.pollInterval, "p", "metrics polling frequency")

	flag.Parse()
	return options
}
