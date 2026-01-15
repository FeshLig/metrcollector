package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/FeshLig/metrcollector/internal/flags"
)

type Options struct {
	Address flags.NetAddress
}

func GetOptions() Options {

	options := Options{
		Address: flags.NetAddress{
			Host: "localhost",
			Port: 8080,
		},
	}

	parseFlags(&options)
	parseEnv(&options)

	return options

}

func parseFlags(options *Options) {

	flag.Var(&options.Address, "a", "net address host:port")

	flag.Parse()

}

func parseEnv(options *Options) error {

	if addrStr, ok := os.LookupEnv("ADDRESS"); ok {
		err := options.Address.Set(addrStr)
		if err != nil {
			return fmt.Errorf("wrong value of ADDRESS: %w", err)
		}
	}

	return nil

}
