package config

import (
	"flag"

	net "github.com/FeshLig/metrcollector/internal/flags"
)

type Options struct {
	Address net.NetAddress
}

func ParseFlags() Options {
	var options Options

	options.Address = net.NetAddress{
		Host: "localhost",
		Port: 8080,
	}

	flag.Var(&options.Address, "a", "net address host:port")

	flag.Parse()
	return options
}
