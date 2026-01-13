package agent

import (
	"flag"
	"time"

	net "github.com/FeshLig/metrcollector/internal/flags"
)

type Options struct {
	Address        net.NetAddress
	ReportInterval net.SecondsDuration
	PollInterval   net.SecondsDuration
}

func ParseFlags() Options {

	var options Options

	const defaultReportInterval = time.Duration(10) * time.Second
	const defaultPollInterval = time.Duration(2) * time.Second

	options.Address = net.NetAddress{
		Host: "localhost",
		Port: 8080,
	}
	options.ReportInterval.Duration = defaultReportInterval
	options.PollInterval.Duration = defaultPollInterval

	flag.Var(&options.Address, "a", "net address host:port")
	flag.Var(&options.ReportInterval, "r", "frequency of sending metrics")
	flag.Var(&options.PollInterval, "p", "metrics polling frequency")

	flag.Parse()
	return options

}
