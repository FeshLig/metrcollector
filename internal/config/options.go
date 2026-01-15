package config

import (
	"flag"

	"github.com/FeshLig/metrcollector/internal/flags"
)

type Options struct {
	Address flags.NetAddress
}

func ParseFlags() Options {
	var options Options

	options.Address = flags.NetAddress{
		Host: "localhost",
		Port: 8080,
	}

	flag.Var(&options.Address, "a", "net address host:port")

	flag.Parse()
	return options
}
