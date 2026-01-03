package main

import (
	"flag"

	"github.com/FeshLig/metrcollector/internal/net"
)

type Options struct {
	address net.NetAddress
}

func ParseFlags() Options {
	var options Options

	options.address = net.NetAddress{
		Host: "localhost",
		Port: 8080,
	}

	flag.Var(&options.address, "a", "net address host:port")

	flag.Parse()
	return options
}
